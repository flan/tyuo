package context
import (
    "time"
)


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
    occurrences int
    lastObserved int64
}

func transitionsIncrement(m map[int]transitionSpec, dictionaryId int) {
    ts, _ := m[dictionaryId] //the nil case for ts will set occurrents to 0
    m[dictionaryId] = transitionSpec{
        occurrences: ts.occurrences + 1,
        lastObserved: time.Now().Unix(),
    }
}
//called before writing the value to the database
func transitionsRescale(m map[int]transitionSpec, rescaleThreshold int,  rescaleDecimator int) {
    rescaleNeeded := false
    for _, ts := range m {
        if ts.occurrences > rescaleThreshold {
            rescaleNeeded = true
            break
        }
    }
    if rescaleNeeded {
        for did, ts := range m {
            ts.occurrences /= rescaleDecimator
            if ts.occurrences > 0 {
                m[did] = ts //it's a copy
            } else {
                delete(m, did)
            }
        }
    }
}

type DigramSpec struct {
    DictionaryIdFirst int
}
type Digram struct {
    transitions map[int]transitionSpec
    
    dictionaryIdFirst int
}
func (g *Digram) rescale(rescaleThreshold int,  rescaleDecimator int) {
    transitionsRescale(g.transitions, rescaleThreshold, rescaleDecimator)
}
func (g *Digram) increment(dictionaryId int) {
    transitionsIncrement(g.transitions, dictionaryId) 
}


type TrigramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
}
type Trigram struct {
    transitions map[int]transitionSpec
    
    dictionaryIdFirst int
    dictionaryIdSecond int
}
func (g *Trigram) rescale(rescaleThreshold int,  rescaleDecimator int) {
    transitionsRescale(g.transitions, rescaleThreshold, rescaleDecimator)
}
func (g *Trigram) increment(dictionaryId int) {
    transitionsIncrement(g.transitions, dictionaryId) 
}

type QuadgramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
}
type Quadgram struct {
    transitions map[int]transitionSpec
    
    dictionaryIdFirst int
    dictionaryIdSecond int
    dictionaryIdThird int
}
func (g *Quadgram) rescale(rescaleThreshold int,  rescaleDecimator int) {
    transitionsRescale(g.transitions, rescaleThreshold, rescaleDecimator)
}
func (g *Quadgram) increment(dictionaryId int) {
    transitionsIncrement(g.transitions, dictionaryId) 
}

type QuintgramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
    DictionaryIdFourth int
}
type Quintgram struct {
    transitions map[int]transitionSpec
    
    dictionaryIdFirst int
    dictionaryIdSecond int
    dictionaryIdThird int
    dictionaryIdFourth int
}
func (g *Quintgram) rescale(rescaleThreshold int,  rescaleDecimator int) {
    transitionsRescale(g.transitions, rescaleThreshold, rescaleDecimator)
}
func (g *Quintgram) increment(dictionaryId int) {
    transitionsIncrement(g.transitions, dictionaryId) 
}




func learnDigramsForward(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specs := make(map[DigramSpec]bool, len(tokens))
    for i := 0; i < len(tokens); i++ {
        specs[DigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
        }] = false
    }
    
    digrams, err := database.digramsGet(specs, true, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    for i := 0; i < len(tokens) - 1; i++ {
        digram := digrams[DigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
        }]
        digram.increment(tokensMap[tokens[i + 1]])
    }
    digram := digrams[DigramSpec{
        DictionaryIdFirst: tokensMap[tokens[len(tokens) - 1]],
    }]
    digram.increment(BoundaryId)
    
    return database.digramsSet(digrams, true, rescaleThreshold, rescaleDecimator)
}
func learnDigramsReverse(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specs := make(map[DigramSpec]bool, len(tokens))
    for i := len(tokens) - 1; i >= 0; i-- {
        specs[DigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
        }] = false
    }
    
    digrams, err := database.digramsGet(specs, false, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    for i := len(tokens) - 1; i >= 1; i-- {
        digram := digrams[DigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
        }]
        digram.increment(tokensMap[tokens[i - 1]])
    }
    digram := digrams[DigramSpec{
        DictionaryIdFirst: tokensMap[tokens[0]],
    }]
    digram.increment(BoundaryId)
    
    return database.digramsSet(digrams, false, rescaleThreshold, rescaleDecimator)
}
func learnDigrams(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    if len(tokens) < 1 {
        return nil
    }
    
    if err := learnDigramsForward(
        database,
        tokens,
        tokensMap,
        oldestAllowedTime,
        rescaleThreshold,
        rescaleDecimator,
    ); err != nil {
        return err
    }
    return learnDigramsReverse(
        database,
        tokens,
        tokensMap,
        oldestAllowedTime,
        rescaleThreshold,
        rescaleDecimator,
    )
}




func learnTrigramsForward(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specs := make(map[TrigramSpec]bool, len(tokens))
    for i := 0; i < len(tokens) - 1; i++ {
        specs[TrigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i + 1]],
        }] = false
    }
    
    trigrams, err := database.trigramsGet(specs, true, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    for i := 0; i < len(tokens) - 2; i++ {
        trigram := trigrams[TrigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i + 1]],
        }]
        trigram.increment(tokensMap[tokens[i + 2]])
    }
    trigram := trigrams[TrigramSpec{
        DictionaryIdFirst: tokensMap[tokens[len(tokens) - 2]],
        DictionaryIdSecond: tokensMap[tokens[len(tokens) - 1]],
    }]
    trigram.increment(BoundaryId)
    
    return database.trigramsSet(trigrams, true, rescaleThreshold, rescaleDecimator)
}
func learnTrigramsReverse(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specs := make(map[TrigramSpec]bool, len(tokens))
    for i := len(tokens) - 1; i >= 1; i-- {
        specs[TrigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i - 1]],
        }] = false
    }
    
    trigrams, err := database.trigramsGet(specs, false, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    for i := len(tokens) - 1; i >= 2; i-- {
        trigram := trigrams[TrigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i - 1]],
        }]
        trigram.increment(tokensMap[tokens[i - 2]])
    }
    trigram := trigrams[TrigramSpec{
        DictionaryIdFirst: tokensMap[tokens[1]],
        DictionaryIdSecond: tokensMap[tokens[0]],
    }]
    trigram.increment(BoundaryId)
    
    return database.trigramsSet(trigrams, false, rescaleThreshold, rescaleDecimator)
}
func learnTrigrams(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    if len(tokens) < 2 {
        return nil
    }
    
    if err := learnTrigramsForward(
        database,
        tokens,
        tokensMap,
        oldestAllowedTime,
        rescaleThreshold,
        rescaleDecimator,
    ); err != nil {
        return err
    }
    return learnTrigramsReverse(
        database,
        tokens,
        tokensMap,
        oldestAllowedTime,
        rescaleThreshold,
        rescaleDecimator,
    )
}



func learnQuadgramsForward(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specs := make(map[QuadgramSpec]bool, len(tokens))
    for i := 0; i < len(tokens) - 2; i++ {
        specs[QuadgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i + 1]],
            DictionaryIdThird: tokensMap[tokens[i + 2]],
        }] = false
    }
    
    quadgrams, err := database.quadgramsGet(specs, true, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    for i := 0; i < len(tokens) - 3; i++ {
        quadgram := quadgrams[QuadgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i + 1]],
            DictionaryIdThird: tokensMap[tokens[i + 2]],
        }]
        quadgram.increment(tokensMap[tokens[i + 3]])
    }
    quadgram := quadgrams[QuadgramSpec{
        DictionaryIdFirst: tokensMap[tokens[len(tokens) - 3]],
        DictionaryIdSecond: tokensMap[tokens[len(tokens) - 2]],
        DictionaryIdThird: tokensMap[tokens[len(tokens) - 1]],
    }]
    quadgram.increment(BoundaryId)
    
    return database.quadgramsSet(quadgrams, true, rescaleThreshold, rescaleDecimator)
}
func learnQuadgramsReverse(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specs := make(map[QuadgramSpec]bool, len(tokens))
    for i := len(tokens) - 1; i >= 2; i-- {
        specs[QuadgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i - 1]],
            DictionaryIdThird: tokensMap[tokens[i - 2]],
        }] = false
    }
    
    quadgrams, err := database.quadgramsGet(specs, false, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    for i := len(tokens) - 1; i >= 3; i-- {
        quadgram := quadgrams[QuadgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i - 1]],
            DictionaryIdThird: tokensMap[tokens[i - 2]],
        }]
        quadgram.increment(tokensMap[tokens[i - 3]])
    }
    quadgram := quadgrams[QuadgramSpec{
        DictionaryIdFirst: tokensMap[tokens[2]],
        DictionaryIdSecond: tokensMap[tokens[1]],
        DictionaryIdThird: tokensMap[tokens[0]],
    }]
    quadgram.increment(BoundaryId)
    
    return database.quadgramsSet(quadgrams, false, rescaleThreshold, rescaleDecimator)
}
func learnQuadgrams(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    if len(tokens) < 3 {
        return nil
    }
    
    if err := learnQuadgramsForward(
        database,
        tokens,
        tokensMap,
        oldestAllowedTime,
        rescaleThreshold,
        rescaleDecimator,
    ); err != nil {
        return err
    }
    return learnQuadgramsReverse(
        database,
        tokens,
        tokensMap,
        oldestAllowedTime,
        rescaleThreshold,
        rescaleDecimator,
    )
}
