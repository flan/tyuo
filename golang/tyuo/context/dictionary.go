package context

//TODO: this should be part of context
const caseSensitiveRepresentationThreshold float64 = 0.1

type dictionaryToken struct {
    id int
    caseInsensitiveOccurrences int
    caseInsensitiveRepresentation string
    capitalisedForms map[string]int
}
//has a function to return the most appropriate representation
    //this function takes the representation threshold as an argument
    //the returned value is a ParsedToken where CaseSensitive is what's expected to be used
    //except where language rules have special handling -- see below.

type ParsedToken struct {
    CaseSensitive string
    CaseInsensitive string
}
//NOTE: when generating the case-insensitive form of a word, the language rules might
//do something like say "if an apostrophe occurs in the middle of this token, its
//case-insensitive form is apostrophe-less", while the with-apostrophe version is considered
//capitalised
//in English, this can probably just be a blanket conversion, with the exception of "it's"
//this should catch "im", "didnt", "thats" and other such things, and eliminate incorrect
//pluralised forms
//when choosing how to present it, if the selected token is identical to its insensitive form
//except for whatever delta the language-rules know how to process, then the CaseSensitive
//value is treated as CaseInsensitive for capitalisation purposes
//basic logic: step through both strings one character at a time, discarding apostrophes
//if, when the end of both are reached, all characters along the way matched, then it's an
//apostrophe variant

type dictionary struct {
    //database reference
    
    //latest ID
    //banned tokens
}
func (d *dictionary) learnTokens(tokens []ParsedToken) (error) {
    //queries the database for all existing insensitive forms, then updates the resulting dictionaryToken instances
    //and creates new ones as needed, writing the result back to the database
    return nil
}
//has functions to take a slice of IDs or tokens and return a corresponding map of dictionaryTokens
//(internally builds map[x]voids to speak with the database)
