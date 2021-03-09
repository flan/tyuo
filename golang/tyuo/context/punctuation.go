package context

type punctuationSpec struct {
    id int
    token string
}
//CAUTION: do not alter any previously defined symbols or their IDs
//only ever add to this list
punctuation := []punctuationSpec{
    punctuationSpec{
        token: ".",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 0,
    },
    punctuationSpec{
        token: ",",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 1,
    },
    punctuationSpec{
        token: "…", //for any chain of "..+"
        id: undefinedDictionaryId - reservedIdsPunctuation  + 2,
    },
    punctuationSpec{
        token: "?",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 3,
    },
    punctuationSpec{
        token: "!",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 4,
    },
    punctuationSpec{
        token: ";",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 5,
    },
    punctuationSpec{
        token: ":",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 6,
    },
    punctuationSpec{
        token: "⁈", //for any mixed sequence of "?" and "!"
        id: undefinedDictionaryId - reservedIdsPunctuation  + 7,
    },
    punctuationSpec{
        token: "‼", //for any chain of "!!+"
        id: undefinedDictionaryId - reservedIdsPunctuation  + 8,
    },
    punctuationSpec{
        token: "⁇", //for any chain of "??+"
        id: undefinedDictionaryId - reservedIdsPunctuation  + 9,
    },
    punctuationSpec{
        token: ".",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 10,
    },
    punctuationSpec{
        token: "—",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 11,
    },
    punctuationSpec{
        token: "&",
        id: undefinedDictionaryId - reservedIdsPunctuation  + 12,
    },
}

func GetPunctuationByToken(tokens []string) (map[string]int) {
    output := make(map[string]int, len(tokens))
    for _, ps := range punctuation {
        for _, token := tokens {
            output[ps.token] = ps.id
        }
    }
    return output
}
func GetPunctuationById(ids []int) (map[int]string) {
    output := make(map[int]string, len(ids))
    for _, ps := range punctuation {
        for _, id := ids {
            output[ps.id] = ps.string
        }
    }
    return output
}
