package logic
import (
    "github.com/flan/tyuo/context"
)

//receives a collection of productions;
//scores and groups any surviving results in descending order of value
func score(ctx *context.Context, productions []production) ([][]production, error) {
    
    //use goroutines liberally
    
    return nil, nil
}


//scoring logic:
//All productions start with a score of 0
//each test adds or removes points (typically 1 or 2) depending on how well the production
//satisfies its requirements:
//each primary keyword is worth one point; each secondary worth one
//failing to meet a minimum target length will deduct a point
//  slightly exceeding the target will award one, but greatly exceeding it will award nothing
//a point will be deducted for every repetition of the same token above two counts

//any production with a positive score is a response candidate and will be formatted
//productions are grouped by score and returned is descending order
