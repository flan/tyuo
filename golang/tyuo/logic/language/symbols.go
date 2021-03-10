//emoticons and maybe emoji
//special cases where an otherwise-unacceptable rune-sequence will be valid

//when lexing, assemble a full non-whitespace-delimited token before looking at what it contains, subject to each rune being in either the language's native space
//or this module

//that can probably be a set that the language assembles in an init()

//once a token is done, check it against a list of known structures here before
//trying to evaluate it as punctuation.



