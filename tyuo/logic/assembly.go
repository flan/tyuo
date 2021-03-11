package logic
import (
    "github.com/flan/tyuo/context"
    //"github.com/flan/tyuo/logic/language"
)

//receives a collection of productions grouped in scored buckets;
//produces a collection of rendered strings, grouped in scored buckets
func assemble(ctx *context.Context, scoredProductions [][]production) ([][]string, error) {
    //iterate over all productions to get a set of IDs
    //retrieve DictionaryTokens for all of those IDs
    
    //pass each production and the DictionaryTokens to the appropriate language's
    //production function
    
    //use goroutines liberally
    
    return nil, nil
}
