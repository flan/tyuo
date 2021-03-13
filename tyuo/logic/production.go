package logic
import (
    "reflect"
    
    "github.com/flan/tyuo/context"
)


func produceDecideStop(ctx *context.Context, pathLen int, minLength int) (bool) {
    if pathLen >= minLength {
        if pathLen >= ctx.GetProductionTargetMinLength() {
            if rng.Float32() < ctx.GetProductionTargetStopProbability() {
                return true
            }
        } else {
            if rng.Float32() < ctx.GetProductionStopProbability() {
                return true
            }
        }
    }
    return false
}

func producePrepareReducedKeytokenIdsSet(keytokenIdsSet map[int]bool, chosenId int) (map[int]bool) {
    newKeyTokenIdsSet := make(map[int]bool, len(keytokenIdsSet) - 1)
    for k, v := range keytokenIdsSet {
        if k != chosenId {
            newKeyTokenIdsSet[k] = v
        }
    }
    return newKeyTokenIdsSet
}

func produceFromNgram(ctx *context.Context, path production, minLength int, keytokenIdsSet map[int]bool, forward bool) ([]production, error) {
    searchBranches := ctx.GetProductionSearchBranchesChildren()
    pathLen := len(path)
    stopConsidered := false
    
    transitionsSelected := false
    transitionIds := make([]int, 0, searchBranches)
    productions := make([]production, 0, 1)
    
    if !transitionsSelected && ctx.AreQuintgramsEnabled() && pathLen >= 4 {
        ngramSpec := context.QuintgramSpec{
            DictionaryIdFirst: path[pathLen - 4],
            DictionaryIdSecond: path[pathLen - 3],
            DictionaryIdThird: path[pathLen - 2],
            DictionaryIdFourth: path[pathLen - 1],
        }
        ngrams, err := ctx.GetQuintgrams(map[context.QuintgramSpec]bool{ngramSpec: false}, forward)
        if err != nil {
            return nil, err
        }
        if len(ngrams) > 0 {
            ngram := ngrams[ngramSpec]
            if ngram.IsTerminal() { //this is a potential ending point
                if !stopConsidered {
                    productions = append(productions, path)
                    if produceDecideStop(ctx, pathLen, minLength) {
                        return productions, nil
                    }
                    stopConsidered = true
                }
            }
            if len(keytokenIdsSet) > 0 {
                if preferredTransitions := ngram.ChooseTransitionIds(keytokenIdsSet, 1); len(preferredTransitions) > 0 {
                    transitionIds = preferredTransitions
                    transitionsSelected = true
                    keytokenIdsSet = producePrepareReducedKeytokenIdsSet(keytokenIdsSet, transitionIds[0])
                }
            }
            transitionIds = append(transitionIds, ngram.SelectTransitionIds(searchBranches - len(transitionIds), ctx.GetIdsBannedStatus)...)
            transitionsSelected = len(transitionIds) >= searchBranches
        }
    }
    
    if !transitionsSelected && ctx.AreQuadgramsEnabled() && pathLen >= 3 {
        ngramSpec := context.QuadgramSpec{
            DictionaryIdFirst: path[pathLen - 3],
            DictionaryIdSecond: path[pathLen - 2],
            DictionaryIdThird: path[pathLen - 1],
        }
        ngrams, err := ctx.GetQuadgrams(map[context.QuadgramSpec]bool{ngramSpec: false}, forward)
        if err != nil {
            return nil, err
        }
        if len(ngrams) > 0 {
            ngram := ngrams[ngramSpec]
            if ngram.IsTerminal() { //this is a potential ending point
                if !stopConsidered {
                    productions = append(productions, path)
                    if produceDecideStop(ctx, pathLen, minLength) {
                        return productions, nil
                    }
                    stopConsidered = true
                }
            }
            if len(keytokenIdsSet) > 0 {
                if preferredTransitions := ngram.ChooseTransitionIds(keytokenIdsSet, 1); len(preferredTransitions) > 0 {
                    transitionIds = preferredTransitions
                    transitionsSelected = true
                    keytokenIdsSet = producePrepareReducedKeytokenIdsSet(keytokenIdsSet, transitionIds[0])
                }
            }
            transitionIds = append(transitionIds, ngram.SelectTransitionIds(searchBranches - len(transitionIds), ctx.GetIdsBannedStatus)...)
            transitionsSelected = len(transitionIds) >= searchBranches
        }
    }
    
    if !transitionsSelected && ctx.AreTrigramsEnabled() && pathLen >= 2 {
        ngramSpec := context.TrigramSpec{
            DictionaryIdFirst: path[pathLen - 2],
            DictionaryIdSecond: path[pathLen - 1],
        }
        ngrams, err := ctx.GetTrigrams(map[context.TrigramSpec]bool{ngramSpec: false}, forward)
        if err != nil {
            return nil, err
        }
        if len(ngrams) > 0 {
            ngram := ngrams[ngramSpec]
            if ngram.IsTerminal() { //this is a potential ending point
                if !stopConsidered {
                    productions = append(productions, path)
                    if produceDecideStop(ctx, pathLen, minLength) {
                        return productions, nil
                    }
                    stopConsidered = true
                }
            }
            if len(keytokenIdsSet) > 0 {
                if preferredTransitions := ngram.ChooseTransitionIds(keytokenIdsSet, 1); len(preferredTransitions) > 0 {
                    transitionIds = preferredTransitions
                    transitionsSelected = true
                    keytokenIdsSet = producePrepareReducedKeytokenIdsSet(keytokenIdsSet, transitionIds[0])
                }
            }
            transitionIds = append(transitionIds, ngram.SelectTransitionIds(searchBranches - len(transitionIds), ctx.GetIdsBannedStatus)...)
            transitionsSelected = len(transitionIds) >= searchBranches
        }
    }
    
    if !transitionsSelected && ctx.AreDigramsEnabled() && pathLen >= 1 {
        ngramSpec := context.DigramSpec{
            DictionaryIdFirst: path[pathLen - 1],
        }
        ngrams, err := ctx.GetDigrams(map[context.DigramSpec]bool{ngramSpec: false}, forward)
        if err != nil {
            return nil, err
        }
        if len(ngrams) > 0 {
            ngram := ngrams[ngramSpec]
            if ngram.IsTerminal() { //this is a potential ending point
                if !stopConsidered {
                    productions = append(productions, path)
                    if produceDecideStop(ctx, pathLen, minLength) {
                        return productions, nil
                    }
                    stopConsidered = true
                }
            }
            if len(keytokenIdsSet) > 0 {
                if preferredTransitions := ngram.ChooseTransitionIds(keytokenIdsSet, 1); len(preferredTransitions) > 0 {
                    transitionIds = preferredTransitions
                    transitionsSelected = true
                    keytokenIdsSet = producePrepareReducedKeytokenIdsSet(keytokenIdsSet, transitionIds[0])
                }
            }
            transitionIds = append(transitionIds, ngram.SelectTransitionIds(searchBranches - len(transitionIds), ctx.GetIdsBannedStatus)...)
            transitionsSelected = len(transitionIds) >= searchBranches
        }
    }
    
    if pathLen < ctx.GetProductionMaxLength() {
        for _, transitionId := range transitionIds {
            newPath := make(production, pathLen + 1)
            copy(newPath, path)
            newPath[pathLen] = transitionId
            if childProductions, err := produceFromNgram(ctx, newPath, minLength, keytokenIdsSet, forward); err == nil {
                if len(productions) > 0 {
                    productions = append(productions, childProductions...)
                }
                break
            } else {
                return nil, err
            }
        }
    }
    return productions, nil
}

func produceFromNgramOrigin(ctx *context.Context, starters <-chan production, minLength int, keytokenIdsSet map[int]bool, forward bool, results chan<- production) {
    for starter := range starters {
        if !forward { //reverse to make the search logic consistent
            for i, j := 0, len(starter) - 1; i < j; i, j = i + 1, j - 1 {
                starter[i], starter[j] = starter[j], starter[i]
            }
        }
        
        if productions, err := produceFromNgram(ctx, starter, minLength, keytokenIdsSet, forward); err == nil {
            for _, production := range productions {
                if !forward { //reverse for consistency
                    for i, j := 0, len(production) - 1; i < j; i, j = i + 1, j - 1 {
                        production[i], production[j] = production[j], production[i]
                    }
                }
                results <- production
            }
        } else {
            logger.Errorf("unable to complete n-gram search: %s", err)
        }
    }
    close(results)
}


func produceStarters(ctx *context.Context, id int, forward bool) ([]production, error) {
    //if an n-gram enumeration turns up a banned option, that's just bad luck; carry on and let the fallback strategies deal with it
    
    searchBranchesRemaining := ctx.GetProductionSearchBranchesInitial()
    searchBranchesBoundaryRemaining := ctx.GetProductionSearchBranchesFromBoundaryInitial()
    productions := make([]production, 0, searchBranchesRemaining)
    
    if ctx.AreQuintgramsEnabled() {
        if searchBranchesRemaining > 0 {
            if ngrams, err := ctx.GetQuintgramsOrigin(id, searchBranchesRemaining, forward); err == nil {
                for _, ngram := range ngrams {
                    if !ctx.AreIdsAllowed([]int{
                        ngram.GetDictionaryIdSecond(),
                        ngram.GetDictionaryIdThird(),
                        ngram.GetDictionaryIdFourth(),
                    }) { //contains a banned value
                        continue
                    }
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    if len(transitionIds) > 0 {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdFirst(),
                            ngram.GetDictionaryIdSecond(),
                            ngram.GetDictionaryIdThird(),
                            ngram.GetDictionaryIdFourth(),
                            transitionIds[0],
                        })
                        searchBranchesRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
        if searchBranchesBoundaryRemaining > 0 {
            if ngrams, err := ctx.GetQuintgramsFromBoundary(id, searchBranchesBoundaryRemaining, forward); err == nil {
                for _, ngram := range ngrams {
                    if !ctx.AreIdsAllowed([]int{
                        ngram.GetDictionaryIdThird(),
                        ngram.GetDictionaryIdFourth(),
                    }) { //contains a banned value
                        continue
                    }
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    if len(transitionIds) > 0 {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdSecond(),
                            ngram.GetDictionaryIdThird(),
                            ngram.GetDictionaryIdFourth(),
                            transitionIds[0],
                        })
                        searchBranchesBoundaryRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
    }
    
    if ctx.AreQuadgramsEnabled() {
        if searchBranchesRemaining > 0 {
            if ngrams, err := ctx.GetQuadgramsOrigin(id, searchBranchesRemaining, forward); err == nil {
                for _, ngram := range ngrams {
                    if !ctx.AreIdsAllowed([]int{
                        ngram.GetDictionaryIdSecond(),
                        ngram.GetDictionaryIdThird(),
                    }) { //contains a banned value
                        continue
                    }
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    if len(transitionIds) > 0 {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdFirst(),
                            ngram.GetDictionaryIdSecond(),
                            ngram.GetDictionaryIdThird(),
                            transitionIds[0],
                        })
                        searchBranchesRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
        if searchBranchesBoundaryRemaining > 0 {
            if ngrams, err := ctx.GetQuadgramsFromBoundary(id, searchBranchesBoundaryRemaining, forward); err == nil {
                for _, ngram := range ngrams {
                    if !ctx.AreIdsAllowed([]int{
                        ngram.GetDictionaryIdThird(),
                    }) { //contains a banned value
                        continue
                    }
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    if len(transitionIds) > 0 {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdSecond(),
                            ngram.GetDictionaryIdThird(),
                            transitionIds[0],
                        })
                        searchBranchesBoundaryRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
    }
    
    if ctx.AreTrigramsEnabled() {
        if searchBranchesRemaining > 0 {
            if ngrams, err := ctx.GetTrigramsOrigin(id, searchBranchesRemaining, forward); err == nil {
                for _, ngram := range ngrams {
                    if !ctx.AreIdsAllowed([]int{
                        ngram.GetDictionaryIdSecond(),
                    }) { //contains a banned value
                        continue
                    }
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    if len(transitionIds) > 0 {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdFirst(),
                            ngram.GetDictionaryIdSecond(),
                            transitionIds[0],
                        })
                        searchBranchesRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
        if searchBranchesBoundaryRemaining > 0 {
            trigramSpec := context.TrigramSpec{DictionaryIdFirst: context.BoundaryId, DictionaryIdSecond: id}
            if ngrams, err := ctx.GetTrigrams(map[context.TrigramSpec]bool{trigramSpec: false}, forward); err == nil {
                if len(ngrams) > 0 {
                    ngram := ngrams[trigramSpec]
                    
                    transitionIds := ngram.SelectTransitionIds(searchBranchesBoundaryRemaining, ctx.GetIdsBannedStatus)
                    for _, transitionId := range transitionIds {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdSecond(),
                            transitionId,
                        })
                        searchBranchesBoundaryRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
    }
    
    if ctx.AreDigramsEnabled() {
        if searchBranchesRemaining > 0 {
            digramSpec := context.DigramSpec{DictionaryIdFirst: id}
            if ngrams, err := ctx.GetDigrams(map[context.DigramSpec]bool{digramSpec: false}, forward); err == nil {
                if len(ngrams) > 0 {
                    ngram := ngrams[digramSpec]
                    
                    transitionIds := ngram.SelectTransitionIds(searchBranchesRemaining, ctx.GetIdsBannedStatus)
                    for _, transitionId := range transitionIds {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdFirst(),
                            transitionId,
                        })
                        searchBranchesRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
        //NOTE: no digrams from boundary, since there are no qualifying search criteria
    }
    
    if !forward { //reverse all productions for consistency
        for _, production := range productions {
            for i, j := 0, len(production) - 1; i < j; i, j = i + 1, j - 1 {
                production[i], production[j] = production[j], production[i]
            }
        }
    }
    
    return productions, nil
}


func produceFromKeytokens(ctx *context.Context, ids []int) ([]production, error) {
    maxInitialProductions := (ctx.GetProductionSearchBranchesInitial() + ctx.GetProductionSearchBranchesFromBoundaryInitial()) * len(ids)
    maxParallelOperations := ctx.GetProductionMaxParallelOperations()
    finishedProductions := make([]production, 0, maxInitialProductions * ctx.GetProductionSearchBranchesChildren() * 2)
    
    
    //do forward entries first to avoid clashing cache-locality with reverse-lookup pages
    queue := make(chan production, maxInitialProductions)
    for _, id := range ids {
        if productions, err := produceStarters(ctx, id, true); err == nil {
            for _, production := range productions {
                queue <- production
            }
        } else {
            return nil, err
        }
    }
    close(queue)
    
    goroutineCount := min(maxParallelOperations, len(queue))
    fragmentSources := make([](chan production), 0, goroutineCount)
    cases := make([]reflect.SelectCase, goroutineCount)
    for i := 0 ; i < goroutineCount; i++ {
        fragmentSource := make(chan production, 1)
        fragmentSources = append(fragmentSources, fragmentSource)
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(fragmentSource)}
        go produceFromNgramOrigin(ctx, queue, 0, nil, true, fragmentSource)
    }
    fragments := make([]production, 0, maxInitialProductions * ctx.GetProductionSearchBranchesChildren())
    remaining := len(cases)
    for remaining > 0 {
        chosen, value, received := reflect.Select(cases)
        if !received { //the channel was closed, so stop watching it
            cases[chosen].Chan = reflect.ValueOf(nil)
            remaining--
            continue
        }
        fragments = append(fragments, value.Interface().(production))
    }
    
    //next, do a reverse-search to finish each production
    queue = make(chan production, len(fragments))
    for _, fragment := range fragments {
        queue <- fragment
    }
    close(queue)
    fragments = nil
    
    goroutineCount = min(maxParallelOperations, len(queue))
    finisherSources := make([](chan production), goroutineCount)
    cases = make([]reflect.SelectCase, goroutineCount)
    for i := 0 ; i < goroutineCount; i++ {
        finisherSource := make(chan production, 1)
        finisherSources = append(finisherSources, finisherSource)
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(finisherSource)}
        go produceFromNgramOrigin(ctx, queue, ctx.GetProductionMinLength(), nil, false, finisherSource)
    }
    remaining = len(cases)
    for remaining > 0 {
        chosen, value, received := reflect.Select(cases)
        if !received { //the channel was closed, so stop watching it
            cases[chosen].Chan = reflect.ValueOf(nil)
            remaining--
            continue
        }
        finishedProductions = append(finishedProductions, value.Interface().(production))
    }
    
    
    //forwards-origin productions are done, so now do the reverse paths
    queue = make(chan production, maxInitialProductions)
    for _, id := range ids {
        if productions, err := produceStarters(ctx, id, false); err == nil {
            for _, production := range productions {
                queue <- production
            }
        } else {
            return nil, err
        }
    }
    close(queue)
    
    goroutineCount = min(maxParallelOperations, len(queue))
    fragmentSources = make([](chan production), 0, goroutineCount)
    cases = make([]reflect.SelectCase, goroutineCount)
    for i := 0 ; i < goroutineCount; i++ {
        fragmentSource := make(chan production, 1)
        fragmentSources = append(fragmentSources, fragmentSource)
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(fragmentSource)}
        go produceFromNgramOrigin(ctx, queue, 0, nil, false, fragmentSource)
    }
    fragments = make([]production, 0, maxInitialProductions * ctx.GetProductionSearchBranchesChildren())
    remaining = len(cases)
    for remaining > 0 {
        chosen, value, received := reflect.Select(cases)
        if !received { //the channel was closed, so stop watching it
            cases[chosen].Chan = reflect.ValueOf(nil)
            remaining--
            continue
        }
        fragments = append(fragments, value.Interface().(production))
    }
    
    //next, do a forward-search to finish each production
    queue = make(chan production, len(fragments))
    for _, fragment := range fragments {
        queue <- fragment
    }
    close(queue)
    fragments = nil
    
    goroutineCount = min(maxParallelOperations, len(queue))
    finisherSources = make([](chan production), goroutineCount)
    cases = make([]reflect.SelectCase, goroutineCount)
    for i := 0 ; i < goroutineCount; i++ {
        finisherSource := make(chan production, 1)
        finisherSources = append(finisherSources, finisherSource)
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(finisherSource)}
        go produceFromNgramOrigin(ctx, queue, ctx.GetProductionMinLength(), nil, true, finisherSource)
    }
    remaining = len(cases)
    for remaining > 0 {
        chosen, value, received := reflect.Select(cases)
        if !received { //the channel was closed, so stop watching it
            cases[chosen].Chan = reflect.ValueOf(nil)
            remaining--
            continue
        }
        finishedProductions = append(finishedProductions, value.Interface().(production))
    }
    
    return finishedProductions, nil
}



func produceTerminalStarters(ctx *context.Context, forward bool) ([]production, error) {
    //if an n-gram enumeration turns up a banned option, that's just bad luck; carry on and let the fallback strategies deal with it
    
    searchBranchesBoundaryRemaining := ctx.GetProductionSearchBranchesFromBoundaryInitial()
    productions := make([]production, 0, searchBranchesBoundaryRemaining)
    
    if ctx.AreQuintgramsEnabled() {
        if searchBranchesBoundaryRemaining > 0 {
            if ngrams, err := ctx.GetQuintgramsOrigin(context.BoundaryId, searchBranchesBoundaryRemaining, forward); err == nil {
                for _, ngram := range ngrams {
                    if !ctx.AreIdsAllowed([]int{
                        ngram.GetDictionaryIdSecond(),
                        ngram.GetDictionaryIdThird(),
                        ngram.GetDictionaryIdFourth(),
                    }) { //contains a banned value
                        continue
                    }
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    if len(transitionIds) > 0 {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdSecond(),
                            ngram.GetDictionaryIdThird(),
                            ngram.GetDictionaryIdFourth(),
                            transitionIds[0],
                        })
                        searchBranchesBoundaryRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
    }
    
    if ctx.AreQuadgramsEnabled() {
        if searchBranchesBoundaryRemaining > 0 {
            if ngrams, err := ctx.GetQuadgramsOrigin(context.BoundaryId, searchBranchesBoundaryRemaining, forward); err == nil {
                for _, ngram := range ngrams {
                    if !ctx.AreIdsAllowed([]int{
                        ngram.GetDictionaryIdSecond(),
                        ngram.GetDictionaryIdThird(),
                    }) { //contains a banned value
                        continue
                    }
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    if len(transitionIds) > 0 {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdSecond(),
                            ngram.GetDictionaryIdThird(),
                            transitionIds[0],
                        })
                        searchBranchesBoundaryRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
    }
    
    if ctx.AreTrigramsEnabled() {
        if searchBranchesBoundaryRemaining > 0 {
            if ngrams, err := ctx.GetTrigramsOrigin(context.BoundaryId, searchBranchesBoundaryRemaining, forward); err == nil {
                for _, ngram := range ngrams {
                    if !ctx.AreIdsAllowed([]int{
                        ngram.GetDictionaryIdSecond(),
                    }) { //contains a banned value
                        continue
                    }
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    if len(transitionIds) > 0 {
                        productions = append(productions, production{
                            ngram.GetDictionaryIdSecond(),
                            transitionIds[0],
                        })
                        searchBranchesBoundaryRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
    }
    
    if ctx.AreDigramsEnabled() {
        if searchBranchesBoundaryRemaining > 0 {
            digramSpec := context.DigramSpec{DictionaryIdFirst: context.BoundaryId}
            if ngrams, err := ctx.GetDigrams(map[context.DigramSpec]bool{digramSpec: false}, forward); err == nil {
                if len(ngrams) > 0 {
                    ngram := ngrams[digramSpec]
                    
                    transitionIds := ngram.SelectTransitionIds(1, ctx.GetIdsBannedStatus)
                    for _, transitionId := range transitionIds {
                        productions = append(productions, production{
                            transitionId,
                        })
                        searchBranchesBoundaryRemaining--
                    }
                }
            } else {
                return nil, err
            }
        }
    }
    
    if !forward { //reverse all productions for consistency
        for _, production := range productions {
            for i, j := 0, len(production) - 1; i < j; i, j = i + 1, j - 1 {
                production[i], production[j] = production[j], production[i]
            }
        }
    }
    
    return productions, nil
}


//picks ID as starting points and produces a slice of productions
func produceFromTerminals(ctx *context.Context, keytokenIds []int, countForward int, countReverse int) ([]production, error) {
    keytokenIdsSet := make(map[int]bool, len(keytokenIds))
    for _, id := range keytokenIds {
        keytokenIdsSet[id] = false
    }
    
    maxInitialProductions := ctx.GetProductionSearchBranchesFromBoundaryInitial()
    maxParallelOperations := ctx.GetProductionMaxParallelOperations()
    finishedProductions := make([]production, 0, maxInitialProductions * ctx.GetProductionSearchBranchesChildren() * 2)
    
    
    //do forward entries first for consistency
    queue := make(chan production, maxInitialProductions)
    if productions, err := produceTerminalStarters(ctx, true); err == nil {
        for _, production := range productions {
            queue <- production
        }
    } else {
        return nil, err
    }
    close(queue)
    
    goroutineCount := min(maxParallelOperations, len(queue))
    finisherSources := make([](chan production), goroutineCount)
    cases := make([]reflect.SelectCase, goroutineCount)
    for i := 0 ; i < goroutineCount; i++ {
        finisherSource := make(chan production, 1)
        finisherSources = append(finisherSources, finisherSource)
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(finisherSource)}
        go produceFromNgramOrigin(ctx, queue, ctx.GetProductionMinLength(), keytokenIdsSet, true, finisherSource)
    }
    remaining := len(cases)
    for remaining > 0 {
        chosen, value, received := reflect.Select(cases)
        if !received { //the channel was closed, so stop watching it
            cases[chosen].Chan = reflect.ValueOf(nil)
            remaining--
            continue
        }
        finishedProductions = append(finishedProductions, value.Interface().(production))
    }
    
    
    //forwards-origin productions are done, so now do the reverse paths
    queue = make(chan production, maxInitialProductions)
    if productions, err := produceTerminalStarters(ctx, false); err == nil {
        for _, production := range productions {
            queue <- production
        }
    } else {
        return nil, err
    }
    close(queue)
    
    goroutineCount = min(maxParallelOperations, len(queue))
    finisherSources = make([](chan production), goroutineCount)
    cases = make([]reflect.SelectCase, goroutineCount)
    for i := 0 ; i < goroutineCount; i++ {
        finisherSource := make(chan production, 1)
        finisherSources = append(finisherSources, finisherSource)
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(finisherSource)}
        go produceFromNgramOrigin(ctx, queue, ctx.GetProductionMinLength(), keytokenIdsSet, false, finisherSource)
    }
    remaining = len(cases)
    for remaining > 0 {
        chosen, value, received := reflect.Select(cases)
        if !received { //the channel was closed, so stop watching it
            cases[chosen].Chan = reflect.ValueOf(nil)
            remaining--
            continue
        }
        finishedProductions = append(finishedProductions, value.Interface().(production))
    }
    
    return finishedProductions, nil
}
