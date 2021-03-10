package context

type punctuationSpec struct {
    id int
    token string
}
//CAUTION: do not alter any previously defined symbols or their IDs;
//only ever add to this list
var punctuation = []punctuationSpec{
    punctuationSpec{
        token: ".",
        id: -2147483647 + 0,
    },
    punctuationSpec{
        token: ",",
        id: -2147483647 + 1,
    },
    punctuationSpec{
        token: "…", //for any chain of "..+"
        id: -2147483647 + 2,
    },
    punctuationSpec{
        token: "?",
        id: -2147483647 + 3,
    },
    punctuationSpec{
        token: "!",
        id: -2147483647 + 4,
    },
    punctuationSpec{
        token: ";",
        id: -2147483647 + 5,
    },
    punctuationSpec{
        token: ":",
        id: -2147483647 + 6,
    },
    punctuationSpec{
        token: "⁈", //for any mixed sequence of "?" and "!"
        id: -2147483647 + 7,
    },
    punctuationSpec{
        token: "‼", //for any chain of "!!+"
        id: -2147483647 + 8,
    },
    punctuationSpec{
        token: "⁇", //for any chain of "??+"
        id: -2147483647 + 9,
    },
    punctuationSpec{
        token: ".",
        id: -2147483647 + 10,
    },
    punctuationSpec{
        token: "—",
        id: -2147483647 + 11,
    },
    punctuationSpec{
        token: "&",
        id: -2147483647 + 12,
    },
} //there's an upper limit of `reservedIdsPunctuation` elements on this structure

//CAUTION: do not write code that alters these structures at runtime
var PunctuationIdsByToken map[string]int = make(map[string]int, len(punctuation))
var PunctuationTokensById map[int]string = make(map[int]string, len(punctuation))
func init() {
    for _, ps := range punctuation {
        PunctuationIdsByToken[ps.token] = ps.id
        PunctuationTokensById[ps.id] = ps.token
    }
}
