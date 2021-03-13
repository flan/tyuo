package logic
import (
    "reflect"
    
    "github.com/flan/tyuo/context"
)





/*
    //how many paths each child should enumerate (but not necessarily explore)
    SearchBranchesChildren int //try 8

    //the minimum number of tokens that need to be produced
    MinLength int
    //the upper limit on how long a production can be
    MaxLength int
    //the likelihood of stopping production, upon finding a terminal,
    //before reaching the target range
    StopProbability float32

    //the minimum desired length of a production
    TargetMinLength int
    //the maximum desired length of a production
    TargetMaxLength int
    //the likelihood of stopping production, upon finding a terminal,
    //after reaching the target range
    TargetStopProbability float32
    //NOTE: for scoring, define "slightly exceeding" as min <= i <= max; "greatly exceeding" as > max
*/
    





//when generating paths from the top level, run each searchBranch in its
//own goroutine, so there should be ten in the base case, all doing reads
//on the database; this should be fine, since only one request can be served
//by each context at any time and creation and learning are separate flows --
//creation is strictly read-only


//when producing, do N forward walks from the keyword and N reverse walks,
//then, for each of the paths that come back (probably grouped by common
//pattern), do a reverse-walk that looks at the full n-gram pattern and
//combine those, rather than the two-start-from-keyword MegaHAL approach.

//if there are no viable chains after scoring, then do N forward and reverse
//walks from the start and end positions, score them, and return that



//the number of sibling-branches to consider for each node in depth-first
//traversal
//the expectation is that each one will produce roughly one result
//Each of the initial branches will be probed, but beyond that, on
//each successive node, each sibling will be tried in order received from the
//database (which is random), and any terminals encountered will produce a
//new forward/reverse option

//NOTE: when looking at transition options, select searchBranches in total, but
//choose anything containing an unencountered keyword first, before doing
//a random/weighted pick of the remainder
//to make this efficient, the keyword-set should be a hash-set, maybe
//with a boolean value set to true/false in each node, if it's expensive to copy it;
//flipped values could be tracked in a slice and restored before returning

//when reaching a stopping point, as the stack is retraversed upwards,
//the next sibling is selected only if there are no candidates from
//traversing the previous branch.

//after falling back to a trigram search, the first terminal found, beyond minDepth,
//will end that discovery path

//when doing a markov walk, choose anything in the keyword set first, if possible

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



func produceFromNgram(ctx *context.Context, path production, keytokenIdsSet map[int]bool, forward bool) ([]production, error) {
    productions := make([]production, 0, 1)
    
    
    //return any productions made in this node and below in the recursive walk
    
    //at each step, if a keytoken is found, copy the set without it and pick that branch;
    //otherwise, just pass the set forward until the walk ends
    
    return productions, nil
}

func produceFromNgramOrigin(ctx *context.Context, starters <-chan production, keytokenIdsSet map[int]bool, forward bool, results chan<- production) {
    //deal with errors in here
    //each instance spawns no goroutines, so respect the stack
    
    //recursively call produceFromNgram, which does the work of querying the database and enumerating transitions
    
    //when it returns, if direction is not forward, reverse each production and join it with the starter before writing it to results
    //otherwise, just prepend the starter and write that
    
    productions := []production{} //returned from produceFromNgram()
    
    for _, production := range productions {
        if !forward { //reverse for consistency
            for i, j := 0, len(production) - 1; i < j; i, j = i + 1, j - 1 {
                production[i], production[j] = production[j], production[i]
            }
        }
        results <- production
    }
}


func produceStarters(ctx *context.Context, id int, forward bool) ([]production, error) {
    //if an n-gram enumeration turns up a banned option, that's just bad luck; carry on and let the fallback strategies deal with it
    
    searchBranchesRemaining := ctx.GetSearchBranchesInitial()
    searchBranchesBoundaryRemaining := ctx.GetSearchBranchesFromBoundaryInitial()
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
    maxInitialProductions := (ctx.GetSearchBranchesInitial() + ctx.GetSearchBranchesFromBoundaryInitial()) * len(ids)
    maxParallelOperations := ctx.GetMaxParallelOperations()
    finishedProductions := make([]production, 0, maxInitialProductions * ctx.GetSearchBranchesChildren() * 2)
    
    
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
    
    fragmentSources := make([](chan production), min(maxParallelOperations, len(queue)))
    cases := make([]reflect.SelectCase, len(fragmentSources))
    for i, fragmentSource := range fragmentSources {
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(fragmentSource)}
        go produceFromNgramOrigin(ctx, queue, nil, true, fragmentSource)
    }
    fragments := make([]production, 0, maxInitialProductions * ctx.GetSearchBranchesChildren())
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
    
    finisherSources := make([](chan production), min(maxParallelOperations, len(queue)))
    cases = make([]reflect.SelectCase, len(finisherSources))
    for i, finisherSource := range finisherSources {
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(finisherSource)}
        go produceFromNgramOrigin(ctx, queue, nil, false, finisherSource)
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
    
    fragmentSources = make([](chan production), min(maxParallelOperations, len(queue)))
    cases = make([]reflect.SelectCase, len(fragmentSources))
    for i, fragmentSource := range fragmentSources {
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(fragmentSource)}
        go produceFromNgramOrigin(ctx, queue, nil, false, fragmentSource)
    }
    fragments = make([]production, 0, maxInitialProductions * ctx.GetSearchBranchesChildren())
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
    
    finisherSources = make([](chan production), min(maxParallelOperations, len(queue)))
    cases = make([]reflect.SelectCase, len(finisherSources))
    for i, finisherSource := range finisherSources {
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(finisherSource)}
        go produceFromNgramOrigin(ctx, queue, nil, true, finisherSource)
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
    
    searchBranchesBoundaryRemaining := ctx.GetSearchBranchesFromBoundaryInitial()
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
    
    maxInitialProductions := ctx.GetSearchBranchesFromBoundaryInitial()
    maxParallelOperations := ctx.GetMaxParallelOperations()
    finishedProductions := make([]production, 0, maxInitialProductions * ctx.GetSearchBranchesChildren() * 2)
    
    
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
    
    finisherSources := make([](chan production), min(maxParallelOperations, len(queue)))
    cases := make([]reflect.SelectCase, len(finisherSources))
    for i, finisherSource := range finisherSources {
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(finisherSource)}
        go produceFromNgramOrigin(ctx, queue, keytokenIdsSet, true, finisherSource)
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
    
    finisherSources = make([](chan production), min(maxParallelOperations, len(queue)))
    cases = make([]reflect.SelectCase, len(finisherSources))
    for i, finisherSource := range finisherSources {
        cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(finisherSource)}
        go produceFromNgramOrigin(ctx, queue, keytokenIdsSet, false, finisherSource)
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
