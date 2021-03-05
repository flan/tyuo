package structure
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("logic")

const sentenceBoundary int = -2147483648 //used to denote the end of a sentence

var rng := rand.Rand.New(rand.NewSource(time.Now().Unix()))
