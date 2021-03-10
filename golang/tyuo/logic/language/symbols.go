package language
import (
    "github.com/flan/tyuo/context"
)

//emoticons and maybe emoji
//special cases where an otherwise-unacceptable rune-sequence will be valid

//when lexing, assemble a full non-whitespace-delimited token before looking at what it contains, subject to each rune being in either the language's native space
//or this module

//that can probably be a set that the language assembles in an init()

//once a token is done, check it against a list of known structures here before
//trying to evaluate it as punctuation.

var symbolRunes = runeset{
    ':': voidInstance,
    ';': voidInstance,
    '<': voidInstance,
    '>': voidInstance,
    '(': voidInstance,
    ')': voidInstance,
    'o': voidInstance,
    '0': voidInstance,
    'O': voidInstance,
    '_': voidInstance,
    '.': voidInstance,
    'T': voidInstance,
    '^': voidInstance,
    'x': voidInstance,
    'n': voidInstance,
    'D': voidInstance,
    '\\': voidInstance,
    '/': voidInstance,
    '3': voidInstance,
    'U': voidInstance,
    'w': voidInstance,
    '🙂': voidInstance,
    '🙁': voidInstance,
}

var symbols = map[string]void{
    ":D": voidInstance,
}

func parseSymbol(token []rune) ([]context.ParsedToken) {
    s := string(token)
    if _, isSymbol := symbols[s]; isSymbol {
        return []context.ParsedToken{
            context.ParsedToken{
                Base: s,
                Variant: s,
            },
        }
    }
    return make([]context.ParsedToken, 0)
}
