package logic
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("logic")

//a forwards-oriented sequence of IDs that describe a produced utterance
type production []int
type scoredProduction struct {
    production production
    score int
    surprise float32
}
type assembledProduction struct {
    Utterance string
    Score int
    Surprise float32
}

var rng = rand.New(rand.NewSource(time.Now().Unix()))
