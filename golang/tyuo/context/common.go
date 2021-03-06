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
const undefinedDictionaryId int = -2147483648

const rescaleThreshold int = 1000
const rescaleDecimator int = 3

var rng = rand.New(rand.NewSource(time.Now().Unix()))
