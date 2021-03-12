package logic
import (
    "github.com/flan/tyuo/context"
)

func produceOrigin(ctx *context.Context, id int, forward bool) ([]production, error) {
    
    //use goroutines sparingly
    
    return nil, nil
}
func produceNgram(ctx *context.Context, path production, forward bool) ([]production, error) {
    
    //use goroutines sparingly
    
    return nil, nil
}






func produceFromNgram(ctx *context.Context, path production, forward bool, results chan<- production) {
    
    
    
}

func produceFromNgramOrigin(ctx *context.Context, starters <-chan production, forward bool, results chan<- production) {
    //deal with errors in here
    //each instance spawns no goroutines, so respect the stack
    
    recursively call produceFromNgram
}


func produceStarters(ctx *context.Context, id int, forward bool) ([]production, error) {
    
    
    
    
    func (c *Context) GetQuintgramsOrigin(
    dictionaryIdFirst int,
    count int,
    forward bool,
) ([]Quintgram, error) {
    return c.database.quintgramsGetOnlyFirst(
        dictionaryIdFirst,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}
    
    
    
    
    func (c *Context) GetQuadgramsOrigin(
    dictionaryIdFirst int,
    count int,
    forward bool,
) ([]Quadgram, error) {
    return c.database.quadgramsGetOnlyFirst(
        dictionaryIdFirst,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}
    
    
    
    
    func (c *Context) GetTrigramsOrigin(
    dictionaryIdFirst int,
    count int,
    forward bool,
) ([]Trigram, error) {
    return c.database.trigramsGetOnlyFirst(
        dictionaryIdFirst,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}


    origin Digrams are just a special case of GetDigrams():
    first = target ID
    if len(results) == 1, then something was found
    

    
}



//picks ID as starting points and produces a slice of productions
func produceFromKeytokens(ctx *context.Context, ids []int) ([]production, error) {
    
    
    
    //generate the initial ngrams based on the selected tokens
    //then create the results channel and kick off a goroutine for each entity,
    //possibly using a workerpool
    
    ctx.GetMaxParallelOperations()
    
    //create a slice of channels for each goroutine and just iterate over that slice,
    //consuming each one until it's closed
    //this will ensure each is able to run to completion and it won't matter if one
    //finishes early
    //when consuming the output, it can be fed into the next channel immediately,
    //even if the goroutines to consume it haven't yet begun
    
    //there'll probably be some sensible way to split the pool, too, once implementation
    //starts
    
    //use goroutines liberally
    
    //for each ID, spawn a bunch of forward and backwards searches
    //if there are fewer IDs that desired for the initial search, pick
    //more origin n-grams to fill out the range
    //for each forward search that produces a result, do a backwards
    //n-gram search; likewise in the other direction
    //maybe this could use a channel so there's a constant pipeline
    
    //it could probably be a workerpool on the channel
    //https://gobyexample.com/worker-pools
    
    return nil, nil
}



func produceTerminalStarters(ctx *context.Context, forward bool) ([]production, error) {
    
    
    func (c *Context) GetQuintgramsFromBoundary(
    dictionaryIdSecond int,
    count int,
    forward bool,
) ([]Quintgram, error) {
    return c.database.quintgramsGetFromBoundary(
        dictionaryIdSecond,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}

    
    func (c *Context) GetQuadgramsFromBoundary(
    dictionaryIdSecond int,
    count int,
    forward bool,
) ([]Quintgram, error) {
    return c.database.quadgramsGetFromBoundary(
        dictionaryIdSecond,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}

    
    boundary trigrams are just a special case of GetTrigrams():
    first = context.BoundaryId
    second = target trigram
    if len(result) == 1 then a match was found
    
    
    boundary digrams are just a special case of GetDigrams():
    first = context.BoundaryId
    if len(result) == 1 then a match was found (which should always happen)
    
}


//picks ID as starting points and produces a slice of productions
func produceFromTerminals(ctx *context.Context, keytokenIds []int, countForward int, countReverse int) ([]production, error) {
    //use goroutines liberally
    
    //query for n-grams in desending order using the terminal ID
    
    //for each ID, do a search in that direction and return whatever comes back
    
    keytokenIdsSet := make(map[int]bool, len(keytokenIds))
    for _, id := range keytokenIds {
        keytokenIdsSet[id] = false
    }
    //at each step, if a keytoken is found, copy the set without it and pick that branch;
    //otherwise, just pass the set forward until the walk ends
    
    return nil, nil
}

/*
    //how many paths to explore from the initial token, in both directions
    SearchBranchesInitial int //try 3
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
