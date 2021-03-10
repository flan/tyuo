package language
import (
    "github.com/flan/tyuo/context"
)

type void struct{}
var voidInstance = void{}

type runeSet map[rune]void

type language struct {
    delimiter rune
    characters runeset
    //TODO: define processing functions
}

var punctuation runeset

func init() {
    punctuation = make(runeset, len(context.PunctuationIdsByToken))
    for punc, _ := range context.PunctuationIdsByToken {
        punctuation[punc] = voidInstance
    }
}
