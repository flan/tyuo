package context
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("context")

type void struct{}
var voidInstance = void{}

//used to denote the end of a sentence, so it will never be a valid ID
const undefinedDictionaryId int32 = -2147483648

const rescaleThreshold int32 = 1000
const rescaleDecimator int32 = 3

var rng = rand.New(rand.NewSource(time.Now().Unix()))
