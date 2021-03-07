package context
import (
    "flag"
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("context")

type void struct{}
var voidInstance = void{}

//used to denote the end of a sentence, so it will never be a valid ID
const undefinedDictionaryId = -2147483648 //int32 minimum; should constrain database sizing

//TODO: these should be flags
const rescaleThreshold = 1000.0
const rescaleDecimator = 3.0
//NOTE: when decimating, use math.RoundToEven to get to 0 faster to slimite obsolete paths

var maxNgramAge = flag.Int64("max-ngram-age", 3600 * 24 * 365, "the number of seconds for which to remember an n-gram value")

var rng = rand.New(rand.NewSource(time.Now().Unix()))
