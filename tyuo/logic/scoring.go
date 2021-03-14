package logic
import (
    "math"
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
    
    var score float32 = 0.0
    
    if len(p) >= ctx.GetProductionMinLength() {
        if len(p) >= ctx.GetProductionTargetMinLength() {
            score += 1.0
        }
    } else {
        score -= 2.0
    }
    
    encounteredTokens := make(map[int]bool, len(p))
    for _, id := range p {
        if _, isKeytoken := keytokenIds[id]; isKeytoken {
            score += 2.0 //award points for keytoken matches
            delete(keytokenIds, id)
        }
        
        if _, alreadyEncountered := encounteredTokens[id]; alreadyEncountered {
            score -= 1.0 //deduct points for repetition
        } else {
            encounteredTokens[id] = false
        }
        
        if _, isPunctuation := context.PunctuationTokensById[id]; isPunctuation {
            score += 0.5 //award points for punctuation, which should offset duplication penalties and favour more interesting phrases
        }
        
        if _, isSymbol := context.SymbolsTokensById[id]; isSymbol {
            score -= 1.0 //remove a point for symbols, making them rarer and dependent on an otherwise-higher-scored production to survive
        }
    }
    
    if score > 0.0 { //if the score isn't positive, don't consider this an option
        output <- scoredProduction{
            production: p,
            score: score,
            surprise: 0.0,
        }
    }
}

func scoreSurpriseQuintgramsGoroutine(
    sp scoredProduction,
    ngrams map[context.QuintgramSpec]context.Quintgram,
    forward bool,
    output chan<- scoredProduction,
    wg *sync.WaitGroup,
) {
    defer wg.Done()
    
    var surprise float32 = 0.0
    production := sp.production
    
    if forward {
        for i := 0; i < len(production) - 3; i++ {
            ngram := ngrams[context.QuintgramSpec{
                DictionaryIdFirst: production[i],
                DictionaryIdSecond: production[i + 1],
                DictionaryIdThird: production[i + 2],
                DictionaryIdFourth: production[i + 3],
            }]
            if i == len(production) - 4 { //terminal position
                surprise += ngram.CalculateSurprise(context.BoundaryId)
            } else {
                surprise += ngram.CalculateSurprise(production[i + 4])
            }
        }
    } else {
        for i := len(production) - 1; i >= 3; i-- {
            ngram := ngrams[context.QuintgramSpec{
                DictionaryIdFirst: production[i],
                DictionaryIdSecond: production[i - 1],
                DictionaryIdThird: production[i - 2],
                DictionaryIdFourth: production[i - 3],
            }]
            if i == 3 { //terminal position
                surprise += ngram.CalculateSurprise(context.BoundaryId)
            } else {
                surprise += ngram.CalculateSurprise(production[i - 4])
            }
        }
    }
    
    sp.surprise = surprise
    output <- sp
}
func scoreSurpriseQuintgrams(ctx *context.Context, scoredProductions []scoredProduction, ngrams map[context.QuintgramSpec]context.Quintgram, forward bool) ([]scoredProduction) {
    var wg sync.WaitGroup
    results := make(chan scoredProduction, len(scoredProductions))
    
    for _, sp := range scoredProductions {
        wg.Add(1)
        go scoreSurpriseQuintgramsGoroutine(sp, ngrams, forward, results, &wg)
    }
    
    surpriseScoredProductions := make([]scoredProduction, 0, len(scoredProductions))
    wg.Wait()
    close(results)
    for sp := range results {
        surpriseScoredProductions = append(surpriseScoredProductions, sp)
    }
    return surpriseScoredProductions
}

func scoreSurpriseQuadgramsGoroutine(
    sp scoredProduction,
    ngrams map[context.QuadgramSpec]context.Quadgram,
    forward bool,
    output chan<- scoredProduction,
    wg *sync.WaitGroup,
) {
    defer wg.Done()
    
    var surprise float32 = 0.0
    production := sp.production
    
    if forward {
        for i := 0; i < len(production) - 2; i++ {
            ngram := ngrams[context.QuadgramSpec{
                DictionaryIdFirst: production[i],
                DictionaryIdSecond: production[i + 1],
                DictionaryIdThird: production[i + 2],
            }]
            if i == len(production) - 3 { //terminal position
                surprise += ngram.CalculateSurprise(context.BoundaryId)
            } else {
                surprise += ngram.CalculateSurprise(production[i + 3])
            }
        }
    } else {
        for i := len(production) - 1; i >= 2; i-- {
            ngram := ngrams[context.QuadgramSpec{
                DictionaryIdFirst: production[i],
                DictionaryIdSecond: production[i - 1],
                DictionaryIdThird: production[i - 2],
            }]
            if i == 2 { //terminal position
                surprise += ngram.CalculateSurprise(context.BoundaryId)
            } else {
                surprise += ngram.CalculateSurprise(production[i - 3])
            }
        }
    }
    
    sp.surprise = surprise
    output <- sp
}
func scoreSurpriseQuadgrams(ctx *context.Context, scoredProductions []scoredProduction, ngrams map[context.QuadgramSpec]context.Quadgram, forward bool) ([]scoredProduction) {
    var wg sync.WaitGroup
    results := make(chan scoredProduction, len(scoredProductions))
    
    for _, sp := range scoredProductions {
        wg.Add(1)
        go scoreSurpriseQuadgramsGoroutine(sp, ngrams, forward, results, &wg)
    }
    
    surpriseScoredProductions := make([]scoredProduction, 0, len(scoredProductions))
    wg.Wait()
    close(results)
    for sp := range results {
        surpriseScoredProductions = append(surpriseScoredProductions, sp)
    }
    return surpriseScoredProductions
}

func scoreSurpriseTrigramsGoroutine(
    sp scoredProduction,
    ngrams map[context.TrigramSpec]context.Trigram,
    forward bool,
    output chan<- scoredProduction,
    wg *sync.WaitGroup,
) {
    defer wg.Done()
    
    var surprise float32 = 0.0
    production := sp.production
    
    if forward {
        for i := 0; i < len(production) - 1; i++ {
            ngram := ngrams[context.TrigramSpec{
                DictionaryIdFirst: production[i],
                DictionaryIdSecond: production[i + 1],
            }]
            if i == len(production) - 2 { //terminal position
                surprise += ngram.CalculateSurprise(context.BoundaryId)
            } else {
                surprise += ngram.CalculateSurprise(production[i + 2])
            }
        }
    } else {
        for i := len(production) - 1; i >= 1; i-- {
            ngram := ngrams[context.TrigramSpec{
                DictionaryIdFirst: production[i],
                DictionaryIdSecond: production[i - 1],
            }]
            if i == 1 { //terminal position
                surprise += ngram.CalculateSurprise(context.BoundaryId)
            } else {
                surprise += ngram.CalculateSurprise(production[i - 2])
            }
        }
    }
    
    sp.surprise = surprise
    output <- sp
}
func scoreSurpriseTrigrams(ctx *context.Context, scoredProductions []scoredProduction, ngrams map[context.TrigramSpec]context.Trigram, forward bool) ([]scoredProduction) {
    var wg sync.WaitGroup
    results := make(chan scoredProduction, len(scoredProductions))
    
    for _, sp := range scoredProductions {
        wg.Add(1)
        go scoreSurpriseTrigramsGoroutine(sp, ngrams, forward, results, &wg)
    }
    
    surpriseScoredProductions := make([]scoredProduction, 0, len(scoredProductions))
    wg.Wait()
    close(results)
    for sp := range results {
        surpriseScoredProductions = append(surpriseScoredProductions, sp)
    }
    return surpriseScoredProductions
}

func scoreSurpriseDigramsGoroutine(
    sp scoredProduction,
    ngrams map[context.DigramSpec]context.Digram,
    forward bool,
    output chan<- scoredProduction,
    wg *sync.WaitGroup,
) {
    defer wg.Done()
    
    var surprise float32 = 0.0
    production := sp.production
    
    if forward {
        for i := 0; i < len(production); i++ {
            ngram := ngrams[context.DigramSpec{
                DictionaryIdFirst: production[i],
            }]
            if i == len(production) - 1 { //terminal position
                surprise += ngram.CalculateSurprise(context.BoundaryId)
            } else {
                surprise += ngram.CalculateSurprise(production[i + 1])
            }
        }
    } else {
        for i := len(production) - 1; i >= 0; i-- {
            ngram := ngrams[context.DigramSpec{
                DictionaryIdFirst: production[i],
            }]
            if i == 0 { //terminal position
                surprise += ngram.CalculateSurprise(context.BoundaryId)
            } else {
                surprise += ngram.CalculateSurprise(production[i - 1])
            }
        }
    }
    
    sp.surprise = surprise
    output <- sp
}
func scoreSurpriseDigrams(ctx *context.Context, scoredProductions []scoredProduction, ngrams map[context.DigramSpec]context.Digram, forward bool) ([]scoredProduction) {
    var wg sync.WaitGroup
    results := make(chan scoredProduction, len(scoredProductions))
    
    for _, sp := range scoredProductions {
        wg.Add(1)
        go scoreSurpriseDigramsGoroutine(sp, ngrams, forward, results, &wg)
    }
    
    surpriseScoredProductions := make([]scoredProduction, 0, len(scoredProductions))
    wg.Wait()
    close(results)
    for sp := range results {
        surpriseScoredProductions = append(surpriseScoredProductions, sp)
    }
    return surpriseScoredProductions
}

func scoreSurprise(ctx *context.Context, scoredProductions []scoredProduction, forward bool) ([]scoredProduction, error) {
    if ctx.AreQuintgramsEnabled() {
        ngramSpecs := make(map[context.QuintgramSpec]bool)
        for _, sp := range scoredProductions {
            production := sp.production
            if forward {
                for i := 0; i < len(production) - 3; i++ {
                    ngramSpecs[context.QuintgramSpec{
                        DictionaryIdFirst: production[i],
                        DictionaryIdSecond: production[i + 1],
                        DictionaryIdThird: production[i + 2],
                        DictionaryIdFourth: production[i + 3],
                    }] = false
                }
            } else {
                for i := len(production) - 1; i >= 3; i-- {
                    ngramSpecs[context.QuintgramSpec{
                        DictionaryIdFirst: production[i],
                        DictionaryIdSecond: production[i - 1],
                        DictionaryIdThird: production[i - 2],
                        DictionaryIdFourth: production[i - 3],
                    }] = false
                }
            }
        }
        ngrams, err := ctx.GetQuintgrams(ngramSpecs, forward)
        if err != nil {
            return nil, err
        }
        
        return scoreSurpriseQuintgrams(ctx, scoredProductions, ngrams, forward), nil
    }
    
    if ctx.AreQuadgramsEnabled() {
        ngramSpecs := make(map[context.QuadgramSpec]bool)
        for _, sp := range scoredProductions {
            production := sp.production
            if forward {
                for i := 0; i < len(production) - 2; i++ {
                    ngramSpecs[context.QuadgramSpec{
                        DictionaryIdFirst: production[i],
                        DictionaryIdSecond: production[i + 1],
                        DictionaryIdThird: production[i + 2],
                    }] = false
                }
            } else {
                for i := len(production) - 1; i >= 2; i-- {
                    ngramSpecs[context.QuadgramSpec{
                        DictionaryIdFirst: production[i],
                        DictionaryIdSecond: production[i - 1],
                        DictionaryIdThird: production[i - 2],
                    }] = false
                }
            }
        }
        ngrams, err := ctx.GetQuadgrams(ngramSpecs, forward)
        if err != nil {
            return nil, err
        }
        
        return scoreSurpriseQuadgrams(ctx, scoredProductions, ngrams, forward), nil
    }
    
    if ctx.AreTrigramsEnabled() {
        ngramSpecs := make(map[context.TrigramSpec]bool)
        for _, sp := range scoredProductions {
            production := sp.production
            if forward {
                for i := 0; i < len(production) - 1; i++ {
                    ngramSpecs[context.TrigramSpec{
                        DictionaryIdFirst: production[i],
                        DictionaryIdSecond: production[i + 1],
                    }] = false
                }
            } else {
                for i := len(production) - 1; i >= 1; i-- {
                    ngramSpecs[context.TrigramSpec{
                        DictionaryIdFirst: production[i],
                        DictionaryIdSecond: production[i - 1],
                    }] = false
                }
            }
        }
        ngrams, err := ctx.GetTrigrams(ngramSpecs, forward)
        if err != nil {
            return nil, err
        }
        
        return scoreSurpriseTrigrams(ctx, scoredProductions, ngrams, forward), nil
    }
    
    if ctx.AreDigramsEnabled() {
        ngramSpecs := make(map[context.DigramSpec]bool)
        for _, sp := range scoredProductions {
            production := sp.production
            if forward {
                for i := 0; i < len(production); i++ {
                    ngramSpecs[context.DigramSpec{
                        DictionaryIdFirst: production[i],
                    }] = false
                }
            } else {
                for i := len(production) - 1; i >= 0; i-- {
                    ngramSpecs[context.DigramSpec{
                        DictionaryIdFirst: production[i],
                    }] = false
                }
            }
        }
        ngrams, err := ctx.GetDigrams(ngramSpecs, forward)
        if err != nil {
            return nil, err
        }
        
        return scoreSurpriseDigrams(ctx, scoredProductions, ngrams, forward), nil
    }
    
    //nothing was done
    return scoredProductions, nil
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
    
    scoredProductions := make([]scoredProduction, 0, len(productions))
    wg.Wait()
    close(results)
    for sp := range results {
        scoredProductions = append(scoredProductions, sp)
    }
    
    
    if ctx.GetProductionCalculateSurpriseForward() || ctx.GetProductionCalculateSurpriseReverse() {
        if ctx.GetProductionCalculateSurpriseForward() {
            sps, err := scoreSurprise(ctx, scoredProductions, true)
            if err != nil {
                return nil, err
            } else {
                scoredProductions = sps
            }
        }
        if ctx.GetProductionCalculateSurpriseReverse() {
            sps, err := scoreSurprise(ctx, scoredProductions, false)
            if err != nil {
                return nil, err
            } else {
                scoredProductions = sps
            }
        }
        
        //numbers and logic stolen directly from MegaHAL, without considering why they were chosen
        for i, sp := range scoredProductions {
            productionLength := len(sp.production)
            surprise := sp.surprise
            if productionLength >= 8 {
                surprise /= float32(math.Sqrt(float64(productionLength - 1)))
            }
            if productionLength >= 16 {
                surprise /= float32(productionLength)
            }
            scoredProductions[i].surprise = surprise
        }
    }
    
    return scoredProductions, nil
}
