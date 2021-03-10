package language
import (
    "github.com/flan/tyuo/context"
)

var symbolRunes = runeset{
    //emoticon and kaomoji bits
    ':': voidInstance,
    ';': voidInstance,
    '<': voidInstance,
    '>': voidInstance,
    '(': voidInstance,
    ')': voidInstance,
    '_': voidInstance,
    '.': voidInstance,
    'T': voidInstance,
    'n': voidInstance,
    'D': voidInstance,
    'o': voidInstance,
    '\\': voidInstance,
    '/': voidInstance,
    '3': voidInstance,
    'c': voidInstance,
    'C': voidInstance,
    '¯': voidInstance,
    'ツ': voidInstance,
    
    //emoji
    '🙂': voidInstance,
    '🙁': voidInstance,
    '🙃': voidInstance,
    '❤️': voidInstance,
}

var symbolsToRepresentation = map[string]string{
    //emoticons
    ":)": ":)",
    ":(": ":(",
    ":D": ":D",
    "D:": "D:",
    ">:D": ">:D",
    "D:<": "D:<",
    ">:(": ">:(",
    "):<": "):<",
    ";_;": ";_;",
    ";.;": ";_;",
    "T_T": "T_T",
    "T.T": "T_T",
    "n.n": "n.n",
    "n_n": "n.n",
    "\o/": "\o/",
    "/o/": "/o/",
    "\\o\\": "\\o\\",
    ":3": ":3",
    ">:3": ">:3",
    "<3": "<3",
    ":C": ":C",
    ":c": ":C",
    ">:C": ">:C",
    ">:c": ">:C",
    ":3c": ":3c",
    
    //kaomoji
    "¯\_(ツ)_/¯": "¯\_(ツ)_/¯",
    
    //emoji
    "🙂": "🙂",
    "🙁": "🙁",
    "🙃": "🙃",
    "❤️": "❤️",
}

func parseSymbol(token []rune) ([]context.ParsedToken) {
    s := string(token)
    if representation, isSymbol := symbolsToRepresentation[s]; isSymbol {
        return []context.ParsedToken{
            context.ParsedToken{
                Base: representation,
                Variant: representation,
            },
        }
    }
    return make([]context.ParsedToken, 0)
}
