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

//picks ID as starting points and produces a slice of productions
func produceFromKeytoken(ctx *context.Context, ids []int) ([]production, error) {
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

//picks ID as starting points and produces a slice of productions
func produceFromTerminals(ctx *context.Context, terminalIdsForward []int, terminalIdsReverse []int) ([]production, error) {
    //use goroutines liberally
    
    //for each ID, do a search in that direction and return whatever comes back
    
    return nil, nil
}

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

