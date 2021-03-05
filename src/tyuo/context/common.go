package structure
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("context")

const rescaleThreshold int = 1000
const rescaleDecimator int = 3

var rng := rand.Rand.New(rand.NewSource(time.Now().Unix()))
