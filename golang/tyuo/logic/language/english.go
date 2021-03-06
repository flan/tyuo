//the idea here is to tokenise input, producing orthogonal and lower-case parallel
//slices, for learning and presenting
//this happens in both flows, before banning is checked... except it's probably
//more performant to just lcase the input and run that through the ban-checker
//then do tokenisation if it's still meaningful to do so
//(which is "always", in the query flow)

//then, for the language-specific bits, in a later pass for the query-stream,
//non-keyword options are filtered out,
//and the most interesting words are selected from the remainder, which means
//there will need to be a reference to the disctionary to select for rarity
