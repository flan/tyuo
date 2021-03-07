package context

//unlike MegaHAL, yuo used trigrams, which will be retained here,
//but it will also introduce quadgrams

//NOTE for use in logic
//when attempting to contruct a sentence, do a quadgram search first,
//changing to trigram on every query-path that doesn't result in a production,
//walking back up the search tree, then proceeding with the next candidate as quadgram

//Each of the initial search-paths is expected to produce roughly one result, so each
//respective search ends when it finds a sentence-boundary token
//a permutation of all combined forward and backwards results is then scored
//and anything over a certain threshold is considered a response candidate;
//all of these are retruned and the requestor can decide which one to present.

//scoring will involve a language-specific component
//in English and French, for example, repeated use of the same token will reduce
//points

//when learning new ngrams, any token in a terminal position gets recorded as a
//terminal in the database
//when producing output, fetch the terminal status of the chosen keyword
//and, if it qualifies, add an empty slice to the forward or backwards glue
//options.

type Terminal struct {
    dictionaryId int
    
    Forward bool
    Reverse bool
}
func (t *Terminal) GetDictionaryId() (int) {
    return t.dictionaryId
}

type transitionSpec struct {
    dictionaryId int
    occurrences int
    lastObserved int
}
type transitions struct {
    transitions map[int]transitionSpec
}
//public function to increment/define transitions

type DigramSpec struct {
    DictionaryIdFirst int
}
type Digram struct {
    transitions
    DigramSpec
}

type TrigramSpec struct {
    DigramSpec
    
    DictionaryIdSecond int
}
type Trigram struct {
    transitions
    TrigramSpec
}

type QuadgramSpec struct {
    TrigramSpec
    
    DictionaryIdThird int
}
type Quadgram struct {
    transitions
    QuadgramSpec
}

type QuintgramSpec struct {
    QuadgramSpec
    
    DictionaryIdFourth int
}
type Quintgram struct {
    transitions
    QuintgramSpec
}
