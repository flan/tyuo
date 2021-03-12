package logic
import (
    "github.com/flan/tyuo/context"
    //"github.com/flan/tyuo/logic/language"
)

//receives a collection of productions with scoring data;
//produces a collection of rendered strings with scoring data
func assemble(ctx *context.Context, scoredProductions []scoredProduction) ([]assembledProduction, error) {
    //iterate over all productions to get a set of IDs
    //retrieve DictionaryTokens for all of those IDs
    
    //pass each production and the DictionaryTokens to the appropriate language's
    //production function
    
    //use goroutines liberally
    
    return nil, nil
}


/*
    //if a token is represented in its base form at least this often,
    //choose that; otherwise, choose the most popular variant
    BaseRepresentationThreshold float32 //0.9 is a good starting point
    */
