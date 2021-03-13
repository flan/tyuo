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
    
    formatUtterance func([]int, map[int]context.DictionaryToken, float32) (string)
}

func getLanguageDefinition(lang string) (*languageDefinition) {
    switch lang {
        case context.LanguageEnglish:
            return &englishLanguageDefinition
        default:
            logger.Errorf("unrecognised language: %s", lang)
            return nil
    }
}
