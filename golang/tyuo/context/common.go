package context
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("context")

type void struct{}
var voidInstance = void{}

//set-types where the value doesn't actually matter
type intSet map[int]bool
type stringSet map[string]bool

//used to denote the end of a sentence
const BoundaryId = -2147483648 //int32 minimum; should constrain database byte-sizing
const undefinedDictionaryId = BoundaryId + 2048 //int32 minimum, plus space for reserved tokens
const reservedIdsPunctuation = 32

var rng = rand.New(rand.NewSource(time.Now().Unix()))
