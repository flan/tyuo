package language
//the idea here is to tokenise input, producing orthogonal and lower-case parallel
//values, for learning and presenting; these are context.ParsedToken structs
//this happens in both flows, before banning is checked... except it's probably
//more performant to just lcase the input and run that through the ban-checker
//then do tokenisation if it's still meaningful to do so
//(which is "always", in the query flow)

//then, for the language-specific bits, in a later pass for the query-stream,
//non-keyword options are filtered out,
//and the most interesting words are selected from the remainder, which means
//there will need to be a reference to the disctionary to select for rarity

//there's also a formatting step where the first token in a sentence gets capitalised,
//if the chosen representation was case-insensitive.


//NOTE: when generating the case-insensitive form of a word, the language rules might
//do something like say "if an apostrophe occurs in the middle of this token, its
//case-insensitive form is apostrophe-less", while the with-apostrophe version is considered
//capitalised
//in English, this can probably just be a blanket conversion, with the exception of "it's"
//this should catch "im", "didnt", "thats" and other such things, and eliminate incorrect
//pluralised forms
//...except this doesn't really work, because there might be valid pluralised forms
//and possessives that would otherwise overlap.
//NOTE: just treat the tokens as they occur; don't try to correct for them
//when choosing how to present it, if the selected token is identical to its insensitive form
//except for whatever delta the language-rules know how to process, then the CaseSensitive
//value is treated as CaseInsensitive for capitalisation purposes
//basic logic: step through both strings one character at a time, discarding apostrophes
//if, when the end of both are reached, all characters along the way matched, then it's an
//apostrophe variant


//when lexing punctuation, convert "--" and standalone "-" into "—"
//convert "!!+" into "‼"
//"??+" into "⁇"
//any chain of "?!" and "!?" into "⁈"
//any sequence "..+" into "…"

var englishLanguage = language{
    delimiter: ' ',
    characters: runeset{
        'a': voidInstance,
        'b': voidInstance,
        'c': voidInstance,
        'd': voidInstance,
        'e': voidInstance,
        'f': voidInstance,
        'g': voidInstance,
        'h': voidInstance,
        'i': voidInstance,
        'j': voidInstance,
        'k': voidInstance,
        'l': voidInstance,
        'm': voidInstance,
        'n': voidInstance,
        'o': voidInstance,
        'p': voidInstance,
        'q': voidInstance,
        'r': voidInstance,
        's': voidInstance,
        't': voidInstance,
        'u': voidInstance,
        'v': voidInstance,
        'w': voidInstance,
        'x': voidInstance,
        'y': voidInstance,
        'z': voidInstance,
        'A': voidInstance,
        'B': voidInstance,
        'C': voidInstance,
        'D': voidInstance,
        'E': voidInstance,
        'F': voidInstance,
        'G': voidInstance,
        'H': voidInstance,
        'I': voidInstance,
        'J': voidInstance,
        'K': voidInstance,
        'L': voidInstance,
        'M': voidInstance,
        'N': voidInstance,
        'O': voidInstance,
        'P': voidInstance,
        'Q': voidInstance,
        'R': voidInstance,
        'S': voidInstance,
        'T': voidInstance,
        'U': voidInstance,
        'V': voidInstance,
        'W': voidInstance,
        'X': voidInstance,
        'Y': voidInstance,
        'Z': voidInstance,
        '\'': voidInstance,
        '-': voidInstance,
    },
}

//when validating whether or not something is a word, if it starts with an apostrophe
//or hyphen, it isn't.
//if it contains any punctuation not in the characterset, it isn't
