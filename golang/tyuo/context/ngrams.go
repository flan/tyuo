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
type Transitions struct {
    transitions map[int]transitionSpec
}
func prepareTransitions(transitions map[int]transitionSpec) (Transitions) {
    return Transitions{
        transitions: transitions,
    }
}
func prepareTransitionsEmpty() (Transitions) {
    return prepareTransitions(make(map[int]transitionSpec, 1))
}
//TODO: should this actually be public? Learning can all happen in-context, and then
//this can just be part of that flow, instead of being a method; see dictionary
//public function to increment/define transitions
func (t *Transitions) Increment(dictionaryId int) {
    ts, _ = t.transitions[dictionaryId] //the nil case for ts will set occurrents to 0
    t.transitions[dictionaryId] = transitionSpec{
        occurrences: ts.occurrences + 1,
        lastObserved: time.Now().Unix(),
    }
}
//called before writing the value to the database
func (t *Transitions) rescale(rescaleThreshold int,  rescaleDeciminator int) {
    recaleNeeded := false
    for _, ts := range t.transitions {
        if ts.occurrences > rescaleThreshold {
            rescaleNeeded = true
            break
        }
    }
    if rescaleNeeded {
        for did, ts := range t.transitions {
            ts.occurrences /= rescaleDeciminator
            if ts.occurrences > 0 {
                t.transitions[did] = ts //it's a copy
            } else {
                delete(t.transitions, did)
            }
        }
    }
}

type DigramSpec struct {
    DictionaryIdFirst int
}
type Digram struct {
    Transitions
    
    dictionaryIdFirst int
}

type TrigramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
}
type Trigram struct {
    Transitions
    
    dictionaryIdFirst int
    dictionaryIdSecond int
}

type QuadgramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
}
type Quadgram struct {
    Transitions
    
    dictionaryIdFirst int
    dictionaryIdSecond int
    dictionaryIdThird int
}

type QuintgramSpec struct {
    DictionaryIdFirst int
    DictionaryIdSecond int
    DictionaryIdThird int
    DictionaryIdFourth int
}
type Quintgram struct {
    Transitions
    
    dictionaryIdFirst int
    dictionaryIdSecond int
    dictionaryIdThird int
    dictionaryIdFourth int
}
