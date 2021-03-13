package context
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
    
    "unicode"
    "golang.org/x/text/cases"
    lang "golang.org/x/text/language"
    "golang.org/x/text/runes"
    "golang.org/x/text/transform"
    "golang.org/x/text/unicode/norm"
)

var logger = loggo.GetLogger("context")

type void struct{}
var voidInstance = void{}

//set-types where the value doesn't actually matter
type intset map[int]bool
type stringset map[string]bool

//used to denote the end of a sentence
const BoundaryId = -2147483648 //int32 minimum; should constrain database byte-sizing
const undefinedDictionaryId = BoundaryId + 4096 //int32 minimum, plus space for reserved tokens
const reservedIdsPunctuation = 32 //from -2147483647 to -2147483615
const reservedIdsPunctuationBase = -2147483647
const reservedIdsSymbols = 1024 //from -2147483614 to -2147482590
const reservedIdsSymbolsBase = -2147483614


var rng = rand.New(rand.NewSource(time.Now().Unix()))


func MakeStringNormaliser() (*transform.Transformer) {
    chain := transform.Chain(
        norm.NFD,
        runes.Remove(runes.In(unicode.Mn)),
        cases.Lower(lang.English), //all normalised data follows English capitalisation rules
        norm.NFC,
    )
    return &chain
}


func intSliceToSet(i []int) (intset) {
    iMap := make(intset, len(i))
    for _, k := range i {
        iMap[k] = false
    }
    return iMap
}
func stringSliceToSet(s []string) (stringset) {
    sMap := make(stringset, len(s))
    for _, k := range s {
        sMap[k] = false
    }
    return sMap
}


func stringSliceToInterfaceSlice(s []string) ([]interface{}) {
    output := make([]interface{}, len(s))
    for i, v := range(s) {
        output[i] = v
    }
    return output
}
func intSliceToInterfaceSlice(s []int) ([]interface{}) {
    output := make([]interface{}, len(s))
    for i, v := range(s) {
        output[i] = v
    }
    return output
}
func stringSetToInterfaceSlice(s stringset) ([]interface{}) {
    output := make([]interface{}, 0, len(s))
    for k, _ := range(s) {
        output = append(output, k)
    }
    return output
}
func intSetToInterfaceSlice(s intset) ([]interface{}) {
    output := make([]interface{}, 0, len(s))
    for k, _ := range(s) {
        output = append(output, k)
    }
    return output
}
