package language
import (
    "github.com/juju/loggo"
    
    "github.com/flan/tyuo/context"
    
    "golang.org/x/text/transform"
)

var logger = loggo.GetLogger("language")

type void struct{}
var voidInstance = void{}

type runeset map[rune]void

type languageDefinition struct {
    delimiter rune
    characters runeset
    
    digestToken func([]rune, *transform.Transformer)([]context.ParsedToken, bool)
}


func Test(x string) {
    //Only here so this package gets compiled
}
