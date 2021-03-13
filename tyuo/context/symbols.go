package context
import (
    "fmt"
)

type symbolSpec struct {
    id int
    repr string
}
func (ss *symbolSpec) GetRepr() (string) {
    return ss.repr
}
func (ss *symbolSpec) GetId() (int) {
    return ss.id
}
//CAUTION: do not alter any previously defined symbols or their IDs;
//only ever add to this list
var Symbols = []symbolSpec{
    //emoticons
    symbolSpec{
        repr: ":)",
        id: reservedIdsSymbolsBase + 0,
    },
    symbolSpec{
        repr: ":(",
        id: reservedIdsSymbolsBase + 1,
    },
    symbolSpec{
        repr: ":|",
        id: reservedIdsSymbolsBase + 2,
    },
    symbolSpec{
        repr: ":D",
        id: reservedIdsSymbolsBase + 3,
    },
    symbolSpec{
        repr: "D:",
        id: reservedIdsSymbolsBase + 4,
    },
    symbolSpec{
        repr: ">:D",
        id: reservedIdsSymbolsBase + 5,
    },
    symbolSpec{
        repr: "D:<",
        id: reservedIdsSymbolsBase + 6,
    },
    symbolSpec{
        repr: ">:(",
        id: reservedIdsSymbolsBase + 7,
    },
    symbolSpec{
        repr: ":O",
        id: reservedIdsSymbolsBase + 8,
    },
    symbolSpec{
        repr: ";_;",
        id: reservedIdsSymbolsBase + 9,
    },
    symbolSpec{
        repr: "T_T",
        id: reservedIdsSymbolsBase + 10,
    },
    symbolSpec{
        repr: "n.n",
        id: reservedIdsSymbolsBase + 11,
    },
    symbolSpec{
        repr: "\\o/",
        id: reservedIdsSymbolsBase + 12,
    },
    symbolSpec{
        repr: "\\o\\",
        id: reservedIdsSymbolsBase + 13,
    },
    symbolSpec{
        repr: "vOv",
        id: reservedIdsSymbolsBase + 14,
    },
    symbolSpec{
        repr: ":3",
        id: reservedIdsSymbolsBase + 15,
    },
    symbolSpec{
        repr: ">:3",
        id: reservedIdsSymbolsBase + 16,
    },
    symbolSpec{
        repr: "<3",
        id: reservedIdsSymbolsBase + 17,
    },
    symbolSpec{
        repr: "</3",
        id: reservedIdsSymbolsBase + 18,
    },
    symbolSpec{
        repr: "C:",
        id: reservedIdsSymbolsBase + 19,
    },
    symbolSpec{
        repr: "C:<",
        id: reservedIdsSymbolsBase + 20,
    },
    symbolSpec{
        repr: ":C",
        id: reservedIdsSymbolsBase + 21,
    },
    symbolSpec{
        repr: ">:C",
        id: reservedIdsSymbolsBase + 22,
    },
    symbolSpec{
        repr: ":3c",
        id: reservedIdsSymbolsBase + 23,
    },
    symbolSpec{
        repr: ">:3c",
        id: reservedIdsSymbolsBase + 24,
    },
    
    
    //kaomoji
    symbolSpec{
        repr: "¯\\_(ツ)_/¯",
        id: reservedIdsSymbolsBase + 127,
    },
    
    
    //emoji
    symbolSpec{
        repr: "😶",
        id: reservedIdsSymbolsBase + 255,
    },
    symbolSpec{
        repr: "😑",
        id: reservedIdsSymbolsBase + 256,
    },
    symbolSpec{
        repr: "🙂",
        id: reservedIdsSymbolsBase + 257,
    },
    symbolSpec{
        repr: "😃",
        id: reservedIdsSymbolsBase + 258,
    },
    symbolSpec{
        repr: "🙁",
        id: reservedIdsSymbolsBase + 259,
    },
    symbolSpec{
        repr: "😦",
        id: reservedIdsSymbolsBase + 260,
    },
    symbolSpec{
        repr: "🙃",
        id: reservedIdsSymbolsBase + 261,
    },
    symbolSpec{
        repr: "🤔",
        id: reservedIdsSymbolsBase + 262,
    },
    symbolSpec{
        repr: "😂",
        id: reservedIdsSymbolsBase + 263,
    },
    symbolSpec{
        repr: "😭",
        id: reservedIdsSymbolsBase + 264,
    },
    symbolSpec{
        repr: "😢",
        id: reservedIdsSymbolsBase + 265,
    },
    symbolSpec{
        repr: "😮",
        id: reservedIdsSymbolsBase + 266,
    },
    symbolSpec{
        repr: "🔥",
        id: reservedIdsSymbolsBase + 267,
    },
} //there's an upper limit of `reservedIdsSymbols` elements on this structure

//a table to ensure consistent representation of symbols
//NOTE: where possible, try to map to similar pre-existing entries that convey the
//same emotional state; it'll help make usage more prevalent and consistent
var symbolsToRepresentation = map[string]string{
    //emoticons
    ":)": ":)",
    "(:": ":)",
    ":(": ":(",
    "):": ":(",
    ":|": ":|",
    "|:": ":|",
    ":D": ":D",
    "D:": "D:",
    ">:D": ">:D",
    "D:<": "D:<",
    ">:(": ">:(",
    "):<": ">:(",
    ":O": ":O",
    ":o": ":O",
    ";_;": ";_;",
    ";.;": ";_;",
    "T_T": "T_T",
    "T.T": "T_T",
    "n.n": "n.n",
    "n_n": "n.n",
    "\\o/": "\\o/",
    "\\o\\": "\\o\\",
    "/o/": "\\o\\",
    "vOv": "vOv",
    ":3": ":3",
    ">:3": ">:3",
    "<3": "<3",
    "</3": "</3",
    "C:": "C:",
    "c:": "C:",
    "C:<": "C:<",
    "c:<": "C:<",
    ":C": ":C",
    ":c": ":C",
    ">:C": ">:C",
    ">:c": ">:C",
    ":3c": ":3c",
    ">:3c": ">:3c",
    
    //kaomoji
    "¯\\_(ツ)_/¯": "¯\\_(ツ)_/¯",
    
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
    "🔥": "🔥",
}

//emoticon and kaomoji bits
var SymbolRunes = make(map[rune]void, len(symbolsToRepresentation))

//CAUTION: do not write code that alters these structures at runtime
var SymbolsIdsByToken map[string]int = make(map[string]int, len(Symbols))
var SymbolsTokensById map[int]string = make(map[int]string, len(Symbols))
func init() {
    if len(Punctuation) > reservedIdsSymbols {
        panic(fmt.Sprintf("symbol-count exceeds reserved limit"))
    }
    
    maxId := reservedIdsSymbolsBase + reservedIdsSymbols
    for _, ss := range Symbols {
        if ss.id < reservedIdsSymbolsBase || ss.id >  + maxId {
            panic(fmt.Sprintf("symbol ID %d is out of range", ss.id))
        }
        
        if _, defined := SymbolsIdsByToken[ss.repr]; defined {
            panic(fmt.Sprintf("duplicate symbol definition for %s", ss.repr))
        }
        SymbolsIdsByToken[ss.repr] = ss.id
        
        if _, defined := SymbolsTokensById[ss.id]; defined {
            panic(fmt.Sprintf("duplicate symbol ID definition for %d", ss.id))
        }
        SymbolsTokensById[ss.id] = ss.repr
    }
    
    //decompose symbols into their parts to enumerate what's allowed in parsing
    for k, _ := range symbolsToRepresentation {
        for _, r := range k {
            SymbolRunes[r] = voidInstance
        }
    }
    
    for _, repr := range symbolsToRepresentation {
        if _, defined := SymbolsIdsByToken[repr]; !defined {
            panic(fmt.Sprintf("symbol %s has no ID entry", repr))
        }
    }
}

func ParseSymbol(token []rune) ([]ParsedToken) {
    s := string(token)
    if representation, isSymbol := symbolsToRepresentation[s]; isSymbol {
        return []ParsedToken{
            ParsedToken{
                Base: representation,
                Variant: representation,
            },
        }
    }
    return nil
}
