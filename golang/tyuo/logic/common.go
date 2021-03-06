package structure
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("logic")

//used to denote the end of a sentence
const sentenceBoundary int = -2147483648

//the number of sibling-branches to consider for each node in depth-first
//traversal
//the expectation is that each one will produce roughly one result
//Each of the initial branches will be probed, but beyond that, on
//each successive node, each sibling will be tried in order received from the
//database (which is random), and any terminals encountered will produce a
//new forward/reverse option
//when reaching a stopping point, as the stack is retraversed upwards,
//the next sibling is selected only if there are no candidates from
//traversing the previous branch.

//after falling back to a trigram search, the first terminal found, beyond minDepth,
//will end that discovery path

//TODO: these should be flag-parameters
const searchBranches int32 = 5
//probability of stopping after each node, upon encountering a terminal
const searchMinDepth int32 = 2
const searchMaxDepth int32 = 30
//probability of stopping on each node, if a terminal is in the set
const searchStopProbability float32 = 0.3

var rng := rand.New(rand.NewSource(time.Now().Unix()))
