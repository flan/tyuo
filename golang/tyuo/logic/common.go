package structure
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("logic")

//used to denote the end of a sentence
const sentenceBoundary int = -2147483648

var rng := rand.New(rand.NewSource(time.Now().Unix()))
