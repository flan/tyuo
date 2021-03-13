package context
import (
    "fmt"
)

type punctuationSpec struct {
    id int
    repr rune
    sentenceInitial bool
}
func (ps *punctuationSpec) GetRepr() (rune) {
    return ps.repr
}
func (ps *punctuationSpec) GetId() (int) {
    return ps.id
}
//CAUTION: do not alter any previously defined punctuation or their IDs;
//only ever add to this list
var Punctuation = []punctuationSpec{
    punctuationSpec{
        repr: '.',
        id: reservedIdsPunctuationBase + 0,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: ',',
        id: reservedIdsPunctuationBase + 1,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: '…', //for any chain of "..+"
        id: reservedIdsPunctuationBase + 2,
        sentenceInitial: true,
    },
    punctuationSpec{
        repr: '?',
        id: reservedIdsPunctuationBase + 3,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: '!',
        id: reservedIdsPunctuationBase + 4,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: ';',
        id: reservedIdsPunctuationBase + 5,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: ':',
        id: reservedIdsPunctuationBase + 6,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: '⁈', //for any mixed sequence of "?" and "!"
        id: reservedIdsPunctuationBase + 7,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: '‼', //for any chain of "!!+"
        id: reservedIdsPunctuationBase + 8,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: '⁇', //for any chain of "??+"
        id: reservedIdsPunctuationBase + 9,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: '—',
        id: reservedIdsPunctuationBase + 10,
        sentenceInitial: false,
    },
    punctuationSpec{
        repr: '&',
        id: reservedIdsPunctuationBase + 11,
        sentenceInitial: false,
    },
} //there's an upper limit of `reservedIdsPunctuation` elements on this structure

//CAUTION: do not write code that alters these structures at runtime
var PunctuationIdsByToken map[string]int = make(map[string]int, len(Punctuation))
var PunctuationTokensById map[int]string = make(map[int]string, len(Punctuation))
var PunctuationTokensNonSentenceInitial map[string]void = make(map[string]void, len(Punctuation) - 1)
func init() {
    if len(Punctuation) > reservedIdsPunctuation {
        panic(fmt.Sprintf("punctuation-count exceeds reserved limit"))
    }
    
    maxId := reservedIdsPunctuationBase + reservedIdsPunctuation
    for _, ps := range Punctuation {
        if ps.id < reservedIdsPunctuationBase || ps.id >  + maxId {
            panic(fmt.Sprintf("punctuation ID %d is out of range", ps.id))
        }
        
        sRepr := string(ps.repr)
        if _, defined := PunctuationIdsByToken[sRepr]; defined {
            panic(fmt.Sprintf("duplicate punctuation definition for %s", sRepr))
        }
        PunctuationIdsByToken[sRepr] = ps.id
        
        if _, defined := SymbolsTokensById[ps.id]; defined {
            panic(fmt.Sprintf("duplicate punctuation ID definition for %d", ps.id))
        }
        PunctuationTokensById[ps.id] = sRepr
        
        if !ps.sentenceInitial {
            PunctuationTokensNonSentenceInitial[sRepr] = voidInstance
        }
    }
}
