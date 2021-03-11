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
    
    var scoredProductions []scoredProduction = nil
    if len(keytokenIds) > 0 {
        //select a subset of the keytokens, to whatever threshold the context wants
        
        productions, err := produceFromKeytokens(ctx, keytokenIds)
        if err != nil {
            logger.Errorf("unable to build productions: %s", err)
            return nil
        }
        scoredProductions, err = score(ctx, productions)
        if err != nil {
            logger.Errorf("unable to score productions: %s", err)
            return nil
        }
    }
    if len(scoredProductions) == 0 { //either no keytokens or no sufficiently good productions
        //select some terminals
        var terminalIdsForward int[] = nil
        var terminalIdsReverse int[] = nil
        
        productions, err := produceFromTerminals(ctx, terminalIdsForward, terminalIdsReverse)
        if err != nil {
            logger.Errorf("unable to build productions: %s", err)
            return nil
        }
        scoredProductions, err = score(ctx, productions)
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
