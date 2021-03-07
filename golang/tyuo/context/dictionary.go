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

type ParsedToken struct {
    CaseSensitive string
    CaseInsensitive string
}

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
