package logic
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("logic")

func max(x, y int) (int) {
    if x < y {
        return y
    }
    return x
}
func min(x, y int) (int) {
    if x > y {
        return y
    }
    return x
}

//a forwards-oriented sequence of IDs that describe a produced utterance
type production []int
type scoredProduction struct {
    production production
    score float32
    surprise float32
}
type assembledProduction struct {
    Utterance string
    Score float32
    Surprise float32
}

var rng = rand.New(rand.NewSource(time.Now().Unix()))
