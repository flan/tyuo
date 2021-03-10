package language
import (
    "github.com/flan/tyuo/context"
)

var symbolRunes = runeset{
    //emoticon and kaomoji bits
    '<': voidInstance,
    '>': voidInstance,
    '(': voidInstance,
    ')': voidInstance,
    '.': voidInstance,
    ':': voidInstance,
    ';': voidInstance,
    '_': voidInstance,
    '\\': voidInstance,
    '/': voidInstance,
    '|': voidInstance.
    '3': voidInstance,
    'C': voidInstance,
    'c': voidInstance,
    'D': voidInstance,
    'n': voidInstance,
    'O': voidInstance,
    'o': voidInstance,
    'T': voidInstance,
    'v': voidInstance,
    '¯': voidInstance,
    'ツ': voidInstance,
    
    //emoji
    '😶': voidInstance,
    '😑': voidInstance,
    '🙂': voidInstance,
    '😃': voidInstance,
    '🙁': voidInstance,
    '😦': voidInstance,
    '🙃': voidInstance,
    '🤔': voidInstance,
    '😂': voidInstance,
    '😭': voidInstance,
    '😢': voidInstance,
    '😮': voidInstance,
    '❤️': voidInstance,
    '💔': voidInstance,
    '🔥': voidInstance,
}

var symbolsToRepresentation = map[string]string{
    //emoticons
    ":)": ":)",
    ":(": ":(",
    ":|": ":|",
    ":D": ":D",
    "D:": "D:",
    ">:D": ">:D",
    "D:<": "D:<",
    ">:(": ">:(",
    "):<": "):<",
    ":O": ":O",
    ":o": ":O",
    ";_;": ";_;",
    ";.;": ";_;",
    "T_T": "T_T",
    "T.T": "T_T",
    "n.n": "n.n",
    "n_n": "n.n",
    "\o/": "\o/",
    "/o/": "/o/",
    "\\o\\": "\\o\\",
    "vOv": "vOv",
    ":3": ":3",
    ">:3": ">:3",
    "<3": "<3",
    "</3": "</3",
    ":C": ":C",
    ":c": ":C",
    ">:C": ">:C",
    ">:c": ">:C",
    ":3c": ":3c",
    ">:3c": ">:3c",
    
    //kaomoji
    "¯\_(ツ)_/¯": "¯\_(ツ)_/¯",
    
    //emoji
    "😶": "😶",
    "😑": "😑",
    "🙂": "🙂",
    "😃": "😃",
    "🙁": "🙁",
    "😦": "😦",
    "🙃": "🙃",
    "🤔": "🤔",
    "😂": "😂",
    "😭": "😭",
    "😢": "😢",
    "😮": "😮",
    "❤️": "❤️",
    "💔": "💔",
    "🔥": "🔥",
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
