package logic
import (
    "sort"
    "sync"
    
    "github.com/flan/tyuo/context"
    "github.com/flan/tyuo/logic/language"
)

func assembleProduction(
    sp scoredProduction,
    dictionaryTokens map[int]context.DictionaryToken,
    output chan<- assembledProduction,
    ctx *context.Context,
    wg *sync.WaitGroup,
) {
    defer wg.Done()
    
    output <- assembledProduction{
        Utterance: language.Format(sp.production, dictionaryTokens, ctx),
        Score: sp.score,
        Surprise: sp.surprise,
    }
}

//receives a collection of productions with scoring data;
//produces a collection of rendered strings with scoring data
func assemble(ctx *context.Context, scoredProductions []scoredProduction) ([]assembledProduction, error) {
    relevantIds := make(map[int]bool)
    for _, sp := range scoredProductions {
        for _, id := range sp.production {
            if _, defined := context.PunctuationTokensById[id]; defined { //filter out punctuation
                continue
            }
            if _, defined := context.SymbolsTokensById[id]; defined { //filter out symbols
                continue
            }
            relevantIds[id] = false
        }
    }
    
    dictionaryTokens, err := ctx.GetDictionaryTokensById(relevantIds)
    if err != nil {
        return nil, err
    }
    
    var wg sync.WaitGroup
    results := make(chan assembledProduction, len(scoredProductions))
    
    for _, sp := range scoredProductions {
        wg.Add(1)
        go assembleProduction(sp, dictionaryTokens, results, ctx, &wg)
    }
    
    assembledProductions := make([]assembledProduction, 0, len(scoredProductions))
    wg.Wait()
    close(results)
    for ap := range results {
        if len(ap.Utterance) > 0 {
            assembledProductions = append(assembledProductions, ap)
        }
    }
    
    sort.Slice(assembledProductions, func(i, j int)(bool){
        if assembledProductions[i].Score == assembledProductions[j].Score {
            return assembledProductions[i].Surprise > assembledProductions[j].Surprise
        }
        return assembledProductions[i].Score > assembledProductions[j].Score
    })
    
    return assembledProductions, nil
}
