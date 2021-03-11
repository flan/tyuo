package logic
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("logic")

//a forwards-oriented sequence of IDs that describe a produced utterance
type production []int

var rng = rand.New(rand.NewSource(time.Now().Unix()))
