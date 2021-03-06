package context
import (
    "math"
    "time"
)


type transitionSpec struct {
    occurrences int
    lastObserved int64
}

func transitionsIncrement(transitions map[int]transitionSpec, dictionaryId int) {
    ts, _ := transitions[dictionaryId] //the nil case for ts will set occurrences to 0
    transitions[dictionaryId] = transitionSpec{
        occurrences: ts.occurrences + 1,
        lastObserved: time.Now().Unix(),
    }
}
//called before writing the value to the database
func transitionsRescale(transitions map[int]transitionSpec, rescaleThreshold int, rescaleDecimator int) {
    rescaleNeeded := false
    for _, ts := range transitions {
        if ts.occurrences > rescaleThreshold {
            rescaleNeeded = true
            break
        }
    }
    if rescaleNeeded {
        for did, ts := range transitions {
            ts.occurrences /= rescaleDecimator
            if ts.occurrences > 0 {
                transitions[did] = ts //it's a copy
            } else {
                delete(transitions, did)
            }
        }
    }
}


func transitionsSumChildren(transitions map[int]transitionSpec) (int) {
    var sum int = 0
    for _, ts := range transitions {
        sum += ts.occurrences
    }
    return sum
}
//this is a weighted random selection of all possible transition nodes,
//a standard Markov-walk selection approach
func transitionsChooseWeightedRandom(
    transitions map[int]transitionSpec,
    count int,
    banCheck func([]int)(map[int]bool),
    excludeBoundaries bool,
) ([]int) {
    remainingTransitions := make(map[int]transitionSpec, len(transitions))
    for did, ts := range transitions {
        if excludeBoundaries && did == BoundaryId {
            continue
        }
        if !banCheck([]int{did})[did] {
            remainingTransitions[did] = ts
        }
    }
    
    selectedIds := make([]int, 0, count)
    for len(selectedIds) < count {
        transitionsSum := transitionsSumChildren(remainingTransitions)
        if transitionsSum == 0 { //all options exhausted
            break
        }
        
        target := rng.Int63n(int64(transitionsSum))
        for dictionaryId, ts := range remainingTransitions {
            target -= int64(ts.occurrences)
            if target <= 0 {
                selectedIds = append(selectedIds, dictionaryId)
                delete(remainingTransitions, dictionaryId)
                break
            }
        }
    }
    return selectedIds
}
func transitionsChooseFromSet(
    transitions map[int]transitionSpec,
    desired map[int]bool,
    count int,
) ([]int) {
    selectedIds := make([]int, 0)
    for k, _ := range desired {
        if _, present := transitions[k]; present {
            selectedIds = append(selectedIds, k)
        }
    }
    if len(selectedIds) > 0 {
        //something was found; randomise what gets picked
        rng.Shuffle(len(selectedIds), func(i, j int) {selectedIds[i], selectedIds[j] = selectedIds[j], selectedIds[i]})
        if len(selectedIds) > count {
            selectedIds = selectedIds[:count]
        }
    }
    return selectedIds
}
//this is part of the surprise-calculation from MegaHAL, used to evaluate how
//predictable a production ended up being as the basis of its scoring system
func transitionsCalculateSurprise(
    transitions map[int]transitionSpec,
    dictionaryId int,
) (float32) {
    ts, defined := transitions[dictionaryId]
    if !defined {
        //an impossible transition was requested, likely because the production
        //n-gram was of a different order
        return 0.0
    }
    
    transitionsSum := transitionsSumChildren(transitions)
    if transitionsSum == 0 {
        //this can happen if an obsolete N-gram is chosen to satisfy the walk's
        //start; just make it neutral
        return 0.0
    }
    
    return float32(-math.Log2(float64(ts.occurrences) / float64(transitionsSum)))
}


type Ngram interface {
    rescale(int, int) 
    increment(int)
    
    IsTerminal() (bool)
    SelectTransitionIds(int, func([]int)(map[int]bool), bool) ([]int)
    ChooseTransitionIds(map[int]bool, int) ([]int)
    CalculateSurprise(int) (float32)
}


type DigramSpec struct {
    DictionaryIdFirst int
}
type Digram struct {
    Ngram
    
    transitions map[int]transitionSpec
    
    dictionaryIdFirst int
}
func (g *Digram) rescale(rescaleThreshold int, rescaleDecimator int) {
    transitionsRescale(g.transitions, rescaleThreshold, rescaleDecimator)
}
func (g *Digram) increment(dictionaryId int) {
    transitionsIncrement(g.transitions, dictionaryId) 
}
func (g *Digram) GetDictionaryIdFirst() (int) {
    return g.dictionaryIdFirst
}
func (g *Digram) IsTerminal() (bool) {
    _, boundary := g.transitions[BoundaryId]
    return boundary
}
func (g *Digram) SelectTransitionIds(
    count int,
    banCheck func([]int)(map[int]bool),
    excludeBoundaries bool,
) ([]int) {
    return transitionsChooseWeightedRandom(g.transitions, count, banCheck, excludeBoundaries)
}
func (g *Digram) ChooseTransitionIds(
    desired map[int]bool,
    count int,
) ([]int) {
    return transitionsChooseFromSet(g.transitions, desired, count)
}
func (g *Digram) CalculateSurprise(dictionaryId int) (float32) {
    return transitionsCalculateSurprise(g.transitions, dictionaryId)
}


type TrigramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
}
type Trigram struct {
    Ngram
    
    transitions map[int]transitionSpec
    
    dictionaryIdFirst int
    dictionaryIdSecond int
}
func (g *Trigram) rescale(rescaleThreshold int, rescaleDecimator int) {
    transitionsRescale(g.transitions, rescaleThreshold, rescaleDecimator)
}
func (g *Trigram) increment(dictionaryId int) {
    transitionsIncrement(g.transitions, dictionaryId) 
}
func (g *Trigram) GetDictionaryIdFirst() (int) {
    return g.dictionaryIdFirst
}
func (g *Trigram) GetDictionaryIdSecond() (int) {
    return g.dictionaryIdSecond
}
func (g *Trigram) IsTerminal() (bool) {
    _, boundary := g.transitions[BoundaryId]
    return boundary
}
func (g *Trigram) SelectTransitionIds(
    count int,
    banCheck func([]int)(map[int]bool),
    excludeBoundaries bool,
) ([]int) {
    return transitionsChooseWeightedRandom(g.transitions, count, banCheck, excludeBoundaries)
}
func (g *Trigram) ChooseTransitionIds(
    desired map[int]bool,
    count int,
) ([]int) {
    return transitionsChooseFromSet(g.transitions, desired, count)
}
func (g *Trigram) CalculateSurprise(dictionaryId int) (float32) {
    return transitionsCalculateSurprise(g.transitions, dictionaryId)
}


type QuadgramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
}
type Quadgram struct {
    Ngram
    
    transitions map[int]transitionSpec
    
    dictionaryIdFirst int
    dictionaryIdSecond int
    dictionaryIdThird int
}
func (g *Quadgram) rescale(rescaleThreshold int, rescaleDecimator int) {
    transitionsRescale(g.transitions, rescaleThreshold, rescaleDecimator)
}
func (g *Quadgram) increment(dictionaryId int) {
    transitionsIncrement(g.transitions, dictionaryId) 
}
func (g *Quadgram) GetDictionaryIdFirst() (int) {
    return g.dictionaryIdFirst
}
func (g *Quadgram) GetDictionaryIdSecond() (int) {
    return g.dictionaryIdSecond
}
func (g *Quadgram) GetDictionaryIdThird() (int) {
    return g.dictionaryIdThird
}
func (g *Quadgram) IsTerminal() (bool) {
    _, boundary := g.transitions[BoundaryId]
    return boundary
}
func (g *Quadgram) SelectTransitionIds(
    count int,
    banCheck func([]int)(map[int]bool),
    excludeBoundaries bool,
) ([]int) {
    return transitionsChooseWeightedRandom(g.transitions, count, banCheck, excludeBoundaries)
}
func (g *Quadgram) ChooseTransitionIds(
    desired map[int]bool,
    count int,
) ([]int) {
    return transitionsChooseFromSet(g.transitions, desired, count)
}
func (g *Quadgram) CalculateSurprise(dictionaryId int) (float32) {
    return transitionsCalculateSurprise(g.transitions, dictionaryId)
}


type QuintgramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
    DictionaryIdFourth int
}
type Quintgram struct {
    Ngram
    
    transitions map[int]transitionSpec
    
    dictionaryIdFirst int
    dictionaryIdSecond int
    dictionaryIdThird int
    dictionaryIdFourth int
}
func (g *Quintgram) rescale(rescaleThreshold int, rescaleDecimator int) {
    transitionsRescale(g.transitions, rescaleThreshold, rescaleDecimator)
}
func (g *Quintgram) increment(dictionaryId int) {
    transitionsIncrement(g.transitions, dictionaryId) 
}
func (g *Quintgram) GetDictionaryIdFirst() (int) {
    return g.dictionaryIdFirst
}
func (g *Quintgram) GetDictionaryIdSecond() (int) {
    return g.dictionaryIdSecond
}
func (g *Quintgram) GetDictionaryIdThird() (int) {
    return g.dictionaryIdThird
}
func (g *Quintgram) GetDictionaryIdFourth() (int) {
    return g.dictionaryIdFourth
}
func (g *Quintgram) IsTerminal() (bool) {
    _, boundary := g.transitions[BoundaryId]
    return boundary
}
func (g *Quintgram) SelectTransitionIds(
    count int,
    banCheck func([]int)(map[int]bool),
    excludeBoundaries bool,
) ([]int) {
    return transitionsChooseWeightedRandom(g.transitions, count, banCheck, excludeBoundaries)
}
func (g *Quintgram) ChooseTransitionIds(
    desired map[int]bool,
    count int,
) ([]int) {
    return transitionsChooseFromSet(g.transitions, desired, count)
}
func (g *Quintgram) CalculateSurprise(dictionaryId int) (float32) {
    return transitionsCalculateSurprise(g.transitions, dictionaryId)
}




func learnDigramsForward(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specOrigin := DigramSpec{
        DictionaryIdFirst: BoundaryId,
    }
    
    specs := make(map[DigramSpec]bool, len(tokens) + 1)
    specs[specOrigin] = false
    for i := 0; i < len(tokens); i++ {
        specs[DigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
        }] = false
    }
    
    digrams, err := database.digramsGet(specs, true, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    digram := digrams[specOrigin]
    digram.increment(BoundaryId)
    
    for i := 0; i < len(tokens) - 1; i++ {
        digram := digrams[DigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
        }]
        digram.increment(tokensMap[tokens[i + 1]])
    }
    
    digram = digrams[DigramSpec{
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
    specOrigin := DigramSpec{
        DictionaryIdFirst: BoundaryId,
    }
    
    specs := make(map[DigramSpec]bool, len(tokens) + 1)
    specs[specOrigin] = false
    for i := len(tokens) - 1; i >= 0; i-- {
        specs[DigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
        }] = false
    }
    
    digrams, err := database.digramsGet(specs, false, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    digram := digrams[specOrigin]
    digram.increment(BoundaryId)
    
    for i := len(tokens) - 1; i >= 1; i-- {
        digram := digrams[DigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
        }]
        digram.increment(tokensMap[tokens[i - 1]])
    }
    
    digram = digrams[DigramSpec{
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
    specOrigin := TrigramSpec{
        DictionaryIdFirst: BoundaryId,
        DictionaryIdSecond: tokensMap[tokens[0]],
    }
    
    specs := make(map[TrigramSpec]bool, len(tokens) + 1)
    specs[specOrigin] = false
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
    
    trigram := trigrams[specOrigin]
    trigram.increment(tokensMap[tokens[1]])
    
    for i := 0; i < len(tokens) - 2; i++ {
        trigram := trigrams[TrigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i + 1]],
        }]
        trigram.increment(tokensMap[tokens[i + 2]])
    }
    
    trigram = trigrams[TrigramSpec{
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
    specOrigin := TrigramSpec{
        DictionaryIdFirst: BoundaryId,
        DictionaryIdSecond: tokensMap[tokens[len(tokens) - 1]],
    }
    
    specs := make(map[TrigramSpec]bool, len(tokens) + 1)
    specs[specOrigin] = false
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
    
    trigram := trigrams[specOrigin]
    trigram.increment(tokensMap[tokens[len(tokens) - 1]])
    
    for i := len(tokens) - 1; i >= 2; i-- {
        trigram := trigrams[TrigramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i - 1]],
        }]
        trigram.increment(tokensMap[tokens[i - 2]])
    }
    
    trigram = trigrams[TrigramSpec{
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
    specOrigin := QuadgramSpec{
        DictionaryIdFirst: BoundaryId,
        DictionaryIdSecond: tokensMap[tokens[0]],
        DictionaryIdThird: tokensMap[tokens[1]],
    }
    
    specs := make(map[QuadgramSpec]bool, len(tokens) + 1)
    specs[specOrigin] = false
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
    
    quadgram := quadgrams[specOrigin]
    quadgram.increment(tokensMap[tokens[2]])
    
    for i := 0; i < len(tokens) - 3; i++ {
        quadgram := quadgrams[QuadgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i + 1]],
            DictionaryIdThird: tokensMap[tokens[i + 2]],
        }]
        quadgram.increment(tokensMap[tokens[i + 3]])
    }
    
    quadgram = quadgrams[QuadgramSpec{
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
    specOrigin := QuadgramSpec{
        DictionaryIdFirst: BoundaryId,
        DictionaryIdSecond: tokensMap[tokens[len(tokens) - 1]],
        DictionaryIdThird: tokensMap[tokens[len(tokens) - 2]],
    }
    
    specs := make(map[QuadgramSpec]bool, len(tokens) + 1)
    specs[specOrigin] = false
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
    
    quadgram := quadgrams[specOrigin]
    quadgram.increment(tokensMap[tokens[len(tokens) - 2]])
    
    for i := len(tokens) - 1; i >= 3; i-- {
        quadgram := quadgrams[QuadgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i - 1]],
            DictionaryIdThird: tokensMap[tokens[i - 2]],
        }]
        quadgram.increment(tokensMap[tokens[i - 3]])
    }
    
    quadgram = quadgrams[QuadgramSpec{
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




func learnQuintgramsForward(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specOrigin := QuintgramSpec{
        DictionaryIdFirst: BoundaryId,
        DictionaryIdSecond: tokensMap[tokens[0]],
        DictionaryIdThird: tokensMap[tokens[1]],
        DictionaryIdFourth: tokensMap[tokens[2]],
    }
    
    specs := make(map[QuintgramSpec]bool, len(tokens) + 1)
    specs[specOrigin] = false
    for i := 0; i < len(tokens) - 3; i++ {
        specs[QuintgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i + 1]],
            DictionaryIdThird: tokensMap[tokens[i + 2]],
            DictionaryIdFourth: tokensMap[tokens[i + 3]],
        }] = false
    }
    
    quintgrams, err := database.quintgramsGet(specs, true, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    quintgram := quintgrams[specOrigin]
    quintgram.increment(tokensMap[tokens[3]])
    
    for i := 0; i < len(tokens) - 4; i++ {
        quintgram := quintgrams[QuintgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i + 1]],
            DictionaryIdThird: tokensMap[tokens[i + 2]],
            DictionaryIdFourth: tokensMap[tokens[i + 3]],
        }]
        quintgram.increment(tokensMap[tokens[i + 4]])
    }
    
    quintgram = quintgrams[QuintgramSpec{
        DictionaryIdFirst: tokensMap[tokens[len(tokens) - 4]],
        DictionaryIdSecond: tokensMap[tokens[len(tokens) - 3]],
        DictionaryIdThird: tokensMap[tokens[len(tokens) - 2]],
        DictionaryIdFourth: tokensMap[tokens[len(tokens) - 1]],
    }]
    quintgram.increment(BoundaryId)
    
    return database.quintgramsSet(quintgrams, true, rescaleThreshold, rescaleDecimator)
}
func learnQuintgramsReverse(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    specOrigin := QuintgramSpec{
        DictionaryIdFirst: BoundaryId,
        DictionaryIdSecond: tokensMap[tokens[len(tokens) - 1]],
        DictionaryIdThird: tokensMap[tokens[len(tokens) - 2]],
        DictionaryIdFourth: tokensMap[tokens[len(tokens) - 3]],
    }
    
    specs := make(map[QuintgramSpec]bool, len(tokens) + 1)
    specs[specOrigin] = false
    for i := len(tokens) - 1; i >= 3; i-- {
        specs[QuintgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i - 1]],
            DictionaryIdThird: tokensMap[tokens[i - 2]],
            DictionaryIdFourth: tokensMap[tokens[i - 3]],
        }] = false
    }
    
    quintgrams, err := database.quintgramsGet(specs, false, oldestAllowedTime)
    if err != nil {
        return err
    }
    
    quintgram := quintgrams[specOrigin]
    quintgram.increment(tokensMap[tokens[len(tokens) - 3]])
    
    for i := len(tokens) - 1; i >= 4; i-- {
        quintgram := quintgrams[QuintgramSpec{
            DictionaryIdFirst: tokensMap[tokens[i]],
            DictionaryIdSecond: tokensMap[tokens[i - 1]],
            DictionaryIdThird: tokensMap[tokens[i - 2]],
            DictionaryIdFourth: tokensMap[tokens[i - 3]],
        }]
        quintgram.increment(tokensMap[tokens[i - 4]])
    }
    
    quintgram = quintgrams[QuintgramSpec{
        DictionaryIdFirst: tokensMap[tokens[3]],
        DictionaryIdSecond: tokensMap[tokens[2]],
        DictionaryIdThird: tokensMap[tokens[1]],
        DictionaryIdFourth: tokensMap[tokens[0]],
    }]
    quintgram.increment(BoundaryId)
    
    return database.quintgramsSet(quintgrams, false, rescaleThreshold, rescaleDecimator)
}
func learnQuintgrams(
    database *database,
    tokens []string,
    tokensMap map[string]int,
    oldestAllowedTime int64,
    rescaleThreshold int,
    rescaleDecimator int,
) (error) {
    if len(tokens) < 4 {
        return nil
    }
    
    if err := learnQuintgramsForward(
        database,
        tokens,
        tokensMap,
        oldestAllowedTime,
        rescaleThreshold,
        rescaleDecimator,
    ); err != nil {
        return err
    }
    return learnQuintgramsReverse(
        database,
        tokens,
        tokensMap,
        oldestAllowedTime,
        rescaleThreshold,
        rescaleDecimator,
    )
}
