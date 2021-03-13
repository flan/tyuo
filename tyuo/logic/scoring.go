package logic
import (
    "github.com/flan/tyuo/context"
)

//receives a collection of productions;
//produces a collection of productions with scoring data
func score(ctx *context.Context, productions []production) ([]scoredProduction, error) {
    scoredProductions := make([]scoredProduction, len(productions))
    
    for _, p := range productions {
        scoredProductions = append(scoredProductions, scoredProduction{
            production: p,
            score: 0,
            surprise: 0.0,
        })
    }
    
    //use goroutines liberally
    
    return scoredProductions, nil
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


//for MegaHAL scoring, use the highest-enabled n-gram level on the context, since that
//should provide decent coverage, reasonable critical value and the lowest amount of I/O;
//in the event that a higher level was enabled post-deployment, there
//will be gaps, but that's fine: just ignore any transitions for which there
//are no records


//if surprise is enabled, also include that value; otherwise, set surprise to 0.0


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
