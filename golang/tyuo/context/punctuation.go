package context

type punctuationSpec struct {
    id int
    repr rune
}
func (ps *punctuationSpec) GetRepr() (rune) {
    return ps.repr
}
func (ps *punctuationSpec) GetId() (int) {
    return ps.id
}
//CAUTION: do not alter any previously defined symbols or their IDs;
//only ever add to this list
var Punctuation = [12]punctuationSpec{
    punctuationSpec{
        repr: '.',
        id: -2147483647 + 0,
    },
    punctuationSpec{
        repr: ',',
        id: -2147483647 + 1,
    },
    punctuationSpec{
        repr: '…', //for any chain of "..+"
        id: -2147483647 + 2,
    },
    punctuationSpec{
        repr: '?',
        id: -2147483647 + 3,
    },
    punctuationSpec{
        repr: '!',
        id: -2147483647 + 4,
    },
    punctuationSpec{
        repr: ';',
        id: -2147483647 + 5,
    },
    punctuationSpec{
        repr: ':',
        id: -2147483647 + 6,
    },
    punctuationSpec{
        repr: '⁈', //for any mixed sequence of "?" and "!"
        id: -2147483647 + 7,
    },
    punctuationSpec{
        repr: '‼', //for any chain of "!!+"
        id: -2147483647 + 8,
    },
    punctuationSpec{
        repr: '⁇', //for any chain of "??+"
        id: -2147483647 + 9,
    },
    punctuationSpec{
        repr: '—',
        id: -2147483647 + 10,
    },
    punctuationSpec{
        repr: '&',
        id: -2147483647 + 11,
    },
} //there's an upper limit of `reservedIdsPunctuation` elements on this structure

//CAUTION: do not write code that alters these structures at runtime
var PunctuationIdsByToken map[string]int = make(map[string]int, len(Punctuation))
var PunctuationTokensById map[int]string = make(map[int]string, len(Punctuation))
func init() {
    for _, ps := range Punctuation {
        PunctuationIdsByToken[string(ps.repr)] = ps.id
        PunctuationTokensById[ps.id] = string(ps.repr)
    }
}
