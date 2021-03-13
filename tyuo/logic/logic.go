package logic
import (
    "runtime/debug"
    
    "github.com/flan/tyuo/context"
    "github.com/flan/tyuo/logic/language"
)

func Speak(ctx *context.Context, input string) ([]assembledProduction) {
    defer func() {
        if r := recover(); r != nil {
            logger.Criticalf(
                "panic observed in Speak(%s): %s\n%s",
                input,
                r,
                string(debug.Stack()),
            )
        }
    }()
    ctx.Lock.RLock()
    defer ctx.Lock.RUnlock()
    
    tokens, _ := language.Parse(input, false, ctx)
    keytokenIds, err := ctx.EnumerateKeytokenIds(tokens)
    if err != nil {
        logger.Errorf("unable to enumerate keytokens: %s", err)
        return nil
    }
    keytokenIdsForScoring := make(map[int]bool, len(keytokenIds))
    for _, id := range keytokenIds {
        keytokenIdsForScoring[id] = false
    }
    
    //number of tokens to start with for each search
    tokensInitial := ctx.GetProductionTokensInitial()
    
    var scoredProductions []scoredProduction = nil
    if len(keytokenIds) > 0 {
        //select a random subset of the keytokens
        rng.Shuffle(len(keytokenIds), func(i, j int){
            keytokenIds[i], keytokenIds[j] = keytokenIds[j], keytokenIds[i]
        })
        if len(keytokenIds) > tokensInitial {
            keytokenIds = keytokenIds[:tokensInitial]
        }
        
        productions, err := produceFromKeytokens(ctx, keytokenIds)
        if err != nil {
            logger.Errorf("unable to build productions: %s", err)
            return nil
        }
        scoredProductions, err = score(ctx, productions, keytokenIdsForScoring)
        if err != nil {
            logger.Errorf("unable to score productions: %s", err)
            return nil
        }
    }
    if len(scoredProductions) == 0 { //either no keytokens or no sufficiently good productions
        countReverse := tokensInitial / 2
        countForward := tokensInitial - countReverse
        //keytokenIds is supplied here, potentially mutated above;
        //if it's not empty, then try to pick them if they come up during the walk;
        //if it is empty, then there's no change to the internal logic
        productions, err := produceFromTerminals(ctx, keytokenIds, countForward, countReverse)
        if err != nil {
            logger.Errorf("unable to build productions: %s", err)
            return nil
        }
        scoredProductions, err = score(ctx, productions, keytokenIdsForScoring)
        if err != nil {
            logger.Errorf("unable to score productions: %s", err)
            return nil
        }
    }
    
    if len(scoredProductions) > 0 {
        assembled, err := assemble(ctx, scoredProductions)
        if err != nil {
            logger.Errorf("unable to assemble productions: %s", err)
            return nil
        }
        return assembled
    }
    return nil
}

func Learn(ctx *context.Context, input []string) (int) {
    defer func() {
        if r := recover(); r != nil {
            logger.Criticalf(
                "panic observed in Learn(%v): %s\n%s",
                input,
                r,
                string(debug.Stack()),
            )
        }
    }()
    ctx.Lock.Lock()
    defer ctx.Lock.Unlock()
    
    linesLearned := 0
    for _, inputLine := range input {
        if !ctx.IsAllowed(inputLine) {
            continue
        }
        
        tokens, learnable := language.Parse(inputLine, true, ctx)
        if learnable && len(tokens) > 0 {
            if err := ctx.LearnInput(tokens); err != nil {
                logger.Errorf("unable to learn input: %s", err)
            } else {
                linesLearned++
            }
        }
    }
    return linesLearned
}

func BanSubstrings(ctx *context.Context, substrings []string) () {
    defer func() {
        if r := recover(); r != nil {
            logger.Criticalf(
                "panic observed in BanSubstrings(%v): %s\n%s",
                substrings,
                r,
                string(debug.Stack()),
            )
        }
    }()
    ctx.Lock.Lock()
    defer ctx.Lock.Unlock()
    
    if err := ctx.BanSubstrings(substrings); err != nil {
        logger.Errorf("unable to ban substrings: %s", err)
    }
}
func UnbanSubstrings(ctx *context.Context, substrings []string) () {
    defer func() {
        if r := recover(); r != nil {
            logger.Criticalf(
                "panic observed in UnbanSubstrings(%v): %s\n%s",
                substrings,
                r,
                string(debug.Stack()),
            )
        }
    }()
    ctx.Lock.Lock()
    defer ctx.Lock.Unlock()
    
    if err := ctx.UnbanSubstrings(substrings); err != nil {
        logger.Errorf("unable to unban substrings: %s", err)
    }
}
