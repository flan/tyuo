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

//used to denote the end of a sentence, so it will never be a valid ID
const undefinedDictionaryId = -2147483648 //int32 minimum; should constrain database sizing

var rng = rand.New(rand.NewSource(time.Now().Unix()))
