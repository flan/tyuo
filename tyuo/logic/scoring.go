package logic
import (
    "sync"
    
    "github.com/flan/tyuo/context"
)


func scoreProduction(
    p production,
    keytokenIds map[int]bool,
    output chan<- scoredProduction,
    ctx *context.Context,
    wg *sync.WaitGroup,
) {
    defer wg.Done()
    
    score := 0
    
    if len(p) >= ctx.GetProductionMinLength() {
        if len(p) >= ctx.GetProductionTargetMinLength() {
            score += 1
        }
    } else {
        score -= 2
    }
    
    encounteredTokens := make(map[int]bool, len(p))
    for _, id := range p {
        if _, isKeytoken := keytokenIds[id]; isKeytoken {
            score += 2 //award points for keytoken matches
            delete(keytokenIds, id)
        }
        
        if _, alreadyEncountered := encounteredTokens[id]; alreadyEncountered {
            score -= 1 //deduct points for repetition
        } else {
            encounteredTokens[id] = false
        }
        
        if _, isPunctuation := context.PunctuationTokensById[id]; isPunctuation {
            score += 1 //award points for punctuation, which should offset duplication penalties and favour more interesting phrases
        }
        
        if _, isSymbol := context.SymbolsTokensById[id]; isSymbol {
            score -= 1 //remove a point for symbols, making them rarer and dependent on an otherwise-higher-scored production to survive
        }
    }
    
    if score > 0 {
        output <- scoredProduction{
            production: p,
            score: score,
            surprise: 0.0,
        }
    }
}

func scoreSurprise() {
    //the goroutine target
}

func scoreSurpriseForward() () {
    
    
    /*
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
    
    */
    
}

func scoreSurpriseReverse() () {
    
}

//receives a collection of productions;
//produces a collection of productions with scoring data
func score(ctx *context.Context, productions []production, keytokenIds map[int]bool) ([]scoredProduction, error) {
    var wg sync.WaitGroup
    results := make(chan scoredProduction, len(productions))
    
    for _, p := range productions {
        wg.Add(1)
        keytokenIdsCopy := make(map[int]bool, len(keytokenIds))
        for k, v := range keytokenIds {
            keytokenIdsCopy[k] = v
        }
        go scoreProduction(p, keytokenIdsCopy, results, ctx, &wg)
    }
    
    scoredProductions := make([]scoredProduction, len(productions))
    wg.Wait()
    close(results)
    for sp := range results {
        scoredProductions = append(scoredProductions, sp)
    }
    
    return scoredProductions, nil
}



/*
func (c *Context) GetDigrams(
    specs map[DigramSpec]bool,
    forward bool,
) (map[DigramSpec]Digram, error) {
    return c.database.digramsGet(
        specs,
        forward,
        c.getOldestAllowedTime(),
    )
}

func (c *Context) GetTrigrams(
    specs map[TrigramSpec]bool,
    forward bool,
) (map[TrigramSpec]Trigram, error) {
    return c.database.trigramsGet(
        specs,
        forward,
        c.getOldestAllowedTime(),
    )
}

func (c *Context) GetQuadgrams(
    specs map[QuadgramSpec]bool,
    forward bool,
) (map[QuadgramSpec]Quadgram, error) {
    return c.database.quadgramsGet(
        specs,
        forward,
        c.getOldestAllowedTime(),
    )
}

func (c *Context) GetQuintgrams(
    specs map[QuintgramSpec]bool,
    forward bool,
) (map[QuintgramSpec]Quintgram, error) {
    return c.database.quintgramsGet(
        specs,
        forward,
        c.getOldestAllowedTime(),
    )
}
*/
