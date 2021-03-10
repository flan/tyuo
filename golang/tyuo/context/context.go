package context
import (
    "sync"
    "time"
)

//context-manager holds a bunch of contexts, keyed by ID
//each context has a database file and a language-specifying file as artifacts
//once a context is loaded, a database connection and the language value are
//held in memory

type contextConfigNgrams struct {
    /* digrams are the simplest and fastest transition model; using them will
     * produce behaviour that is often novel, sometimes insightful,
     * frequently deranged, particularly as learning progresses.
     * It's pretty random and will only resemble native speech by chance.
     */
    Digrams bool
    /* trigrams are a fairly middle-ground option, producing relevant
     * observations with some regularity and having about as much
     * sentence-structure correctness as a machine-translation between
     * languages with no common ancestry.
     */
    Trigrams bool
    /*
     * quadgrams are a reasonably stable choice for production of "how do you
     * do, fellow humans" responses, being sell-formed for the most part, but
     * closely reflecting observed input: a lot of data will need to be
     * learned before novel structures will be produced with any regularity
     * and search-spaces will sometimes be exhausted while finding viable
     * paths.
     */
    Quadgrams bool
    /* quintgrams (and anything above them) will rarely deviate from mimicing
     * what was learned; occasional novel productions are possible, but it
     * will not be uncommon to see near-verbatim recreations of input data.
     */
    Quintgrams bool
}

type contextConfigLearning struct {
    //how long, in tokens, input needs to be before learning will occur;
    //it is automatically fed to any enabled n-gram structures that
    //can accomodate the given length
    MinLength int
    
    //the number of runes allowed within any single token,
    //used to prevent over-hyphenated compounds that will only
    //ever be seen a handful of times from cluttering the database
    MaxTokenLength int //12-15 is probably a good range for this
    
    //how long to hold on to n-gram structures
    MaxAge int64
    
    //the number of dictionary occurrences or transitions at which
    //to trigger rescale logic, which eliminates obsolete entries and
    //keeps the numbers in check
    RescaleThreshold int //should probably be 1000
    //the divisor for rescaling; this affects how frequently it happens
    //and how long rare entries hang around
    RescaleDecimator int //should probably be 3
}

type contextConfigProduction struct {
    //how many paths to explore from the initial token, in both directions
    SearchBranchesInitial int //try 4
    //how many paths each child should enumerate (but not necessarily explore)
    SearchBranchesChildren int //try 10
    
    //the minimum number of tokens that need to be produced
    MinLength int
    //the upper limit on how long a production can be
    MaxLength int
    //the likelihood of stopping production, upon finding a terminal,
    //before reaching the target range
    StopProbability float32
    
    //the minimum desired length of a production
    TargetMinLength int
    //the maximum desired length of a production
    TargetMaxLength int
    //the likelihood of stopping production, upon finding a terminal,
    //after reaching the target range
    TargetStopProbability float32
    //NOTE: for scoring, define "slightly exceeding" as min <= i <= max; "greatly exceeding" as > max
    
    //if a token is represented in its base form at least this often,
    //choose that; otherwise, choose the most popular variant
    BaseRepresentationThreshold float32 //0.9 is a good starting point
}

type contextConfig struct {
    //"en", "fr"
    Language string
    
    Ngrams contextConfigNgrams
    
    Learning contextConfigLearning
    
    Production contextConfigProduction
}

/*
~/.tyuo/
    contexts/
        <id>.sqlite3
        <id>.json {
            "language": "en",
            "nGrams": {
                "digrams": false,
                "trigrams": true,
                "quadgrams": true,
                "quintgrams": true,
            },
            "learning": {
                "minLength": 5, //if learnable input is at least this long, update the dictionary and feed it into any n-grams where it fits
                "maxAge": <a year, in seconds>,
            },
            "production": {
                "minLength": 4,
                "maxLength": 30, //just abort the production search at this point
                "stopProbability": 0.25, //applies after min until reaching the target range
                "targetLength": {
                    "min": 10,
                    "max": 20,
                    "stopProbability": 0.5, //applies until production ends
                    //for scoring, define "slightly exceeding" as min <= i <= max; "greatly exceeding" as > max
                },
                "caseSensitivityThreshold": 0.1,
            },
        }
    languages/
        <language>.banned
        <language>.boring
*/

func intSliceToSet(i []int) (intset) {
    iMap := make(intset, len(i))
    for _, k := range i {
        iMap[k] = false
    }
    return iMap
}
func stringSliceToSet(s []string) (stringset) {
    sMap := make(stringset, len(s))
    for _, k := range s {
        sMap[k] = false
    }
    return sMap
}

type Context struct {
    config contextConfig
    
    database *database
    bannedDictionary *bannedDictionary
    dictionary *dictionary
    boringTokens map[string]void
    
    //users of this struct are expected to respect this lock
    //learning is a writing flow; everything else is reading
    Lock sync.RWMutex
}
func (c *Context) GetLanguage() (string) {
    return c.config.Language
}
func (c *Context) GetMaxTokenLength() (int) {
    return c.config.Learning.MaxTokenLength
}

func (c *Context) AreDigramsEnabled() (bool) {
    return c.config.Ngrams.Digrams
}
func (c *Context) AreTrigramsEnabled() (bool) {
    return c.config.Ngrams.Trigrams
}
func (c *Context) AreQuadgramsEnabled() (bool) {
    return c.config.Ngrams.Quadgrams
}
func (c *Context) AreQuintgramsEnabled() (bool) {
    return c.config.Ngrams.Quintgrams
}




func (c *Context) getOldestAllowedTime() (int64) {
    return time.Now().Unix() - c.config.Learning.MaxAge
}

func (c *Context) GetTerminals(ids []int) (map[int]Terminal, error) {
    return c.database.terminalsGetTerminals(
        intSliceToSet(ids),
        c.getOldestAllowedTime(),
    )
}
//TODO: define other paths for n-gram database access


func (c *Context) IsAllowed(s string) (bool) {
    return !c.bannedDictionary.containsBannedToken(s)
}
func (c *Context) GetIdBannedStatus(ids []int) (map[int]bool) {
    return c.bannedDictionary.getIdBannedStatus(intSliceToSet(ids))
}



func (c *Context) LearnInput(tokens []ParsedToken) (error) {
    if len(tokens) == 0 {
        return nil
    }
    
    rescaleThreshold := c.config.Learning.RescaleThreshold
    rescaleDecimator := c.config.Learning.RescaleDecimator
    
    //strip punctuation so it doesn't get learned redundantly
    depunctuatedTokens := make([]ParsedToken, 0, len(tokens))
    for _, pt := range tokens {
        if _, defined := PunctuationIdsByToken[pt.Base]; !defined {
            depunctuatedTokens = append(depunctuatedTokens, pt)
        }
    }
    
    //first, update the dictionary to make sure all tokens have an ID
    dictionaryTokens, err := c.dictionary.learnTokens(
        depunctuatedTokens,
        rescaleThreshold,
        rescaleDecimator,
    )
    if err != nil {
        return err
    }
    //not needed anymore and this function is far from over
    depunctuatedTokens = nil
    
    tokensMap := make(map[string]int, len(dictionaryTokens) + len(PunctuationIdsByToken))
    for _, dt := range dictionaryTokens {
        tokensMap[dt.baseRepresentation] = dt.id
    }
    //put punctuation mappings in
    for token, id := range PunctuationIdsByToken {
        tokensMap[token] = id
    }
    
    baseTokens := make([]string, len(tokens))
    for i, token := range tokens {
        baseTokens[i] = token.Base
    }
    
    if err = learnTerminals(
        c.database,
        tokensMap[baseTokens[0]],
        tokensMap[baseTokens[len(baseTokens) - 1]],
    ); err != nil {
        return err
    }
    
    oldestAllowedTime := c.getOldestAllowedTime()
    if c.AreDigramsEnabled() {
        if err = learnDigrams(
            c.database,
            baseTokens,
            tokensMap,
            oldestAllowedTime,
            rescaleThreshold,
            rescaleDecimator,
        ); err != nil {
            return err
        }
    }
    if c.AreTrigramsEnabled() && len(tokens) > 1 {
        if err = learnTrigrams(
            c.database,
            baseTokens,
            tokensMap,
            oldestAllowedTime,
            rescaleThreshold,
            rescaleDecimator,
        ); err != nil {
            return err
        }
    }
    if c.AreQuadgramsEnabled() && len(tokens) > 2 {
        if err = learnQuadgrams(
            c.database,
            baseTokens,
            tokensMap,
            oldestAllowedTime,
            rescaleThreshold,
            rescaleDecimator,
        ); err != nil {
            return err
        }
    }
    if c.AreQuintgramsEnabled()  && len(tokens) > 3 {
        if err = learnQuintgrams(
            c.database,
            baseTokens,
            tokensMap,
            oldestAllowedTime,
            rescaleThreshold,
            rescaleDecimator,
        ); err != nil {
            return err
        }
    }
    
    return nil
}

func (c *Context) EnumerateKeytokenIds(tokens []ParsedToken) ([]int, error) {
    candidates := make(stringset, len(tokens))
    for _, pt := range tokens {
        if _, isPunctuation := PunctuationIdsByToken[pt.Base]; isPunctuation {
            continue
        }
        if _, isBoring := c.boringTokens[pt.Base]; isBoring {
            continue
        }
        if c.bannedDictionary.containsBannedToken(pt.Base) {
            continue
        }
        
        candidates[pt.Base] = false
    }
    
    return c.dictionary.getIdsByToken(candidates)
}





type ContextManager struct {
    databaseManager databaseManager
    
    bannedTokensGenericByLanguage map[string][]string
    boringTokensByLanguage map[string]map[string]void
    
    contexts map[string]Context
    
    //used internally to control access to GetContext(), so that
    //resources like the database aren't connected multiple times
    lock sync.Mutex
}
func (cm *ContextManager) Close() {
    cm.databaseManager.Close()
    cm.contexts = make(map[string]Context)
}
func (cm *ContextManager) GetContext(contextId string) (*Context, error) {
    cm.lock.Lock()
    defer cm.lock.Unlock()
    
    return nil, nil
}

//function to get contexts
//nothing except "GetContext()"; contexts need to be defined
//by creating a JSON file on disk.


//Only used to cause compilation of this package
func Test(s string){
}
