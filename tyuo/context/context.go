package context
import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
)

const LanguageEnglish = "english"
const LanguageFrench = "french"


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
    MinTokenCount int //10

    //the number of runes allowed within any single token,
    //used to prevent over-hyphenated compounds that will only
    //ever be seen a handful of times from cluttering the database
    MaxTokenLength int //12-15 is probably a good range for this

    //how long to hold on to n-gram structures, in seconds
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
    //the maximum number of operations to run in parallel at any stage in the
    //production pipeline
    MaxParallelOperations int
    
    //the number of keytokens or terminals to choose before starting a search
    TokensInitial int //try 2
    //how many paths to explore from the initial token, in both directions
    SearchBranchesInitial int //try 3
    //how many paths to explore from the bounary, in both directions
    SearchBranchesFromBoundaryInitial //try 1
    //how many paths each child should enumerate (but not necessarily explore)
    SearchBranchesChildren int //try 8

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
    
    //enabling this calculates surprise for each production, similar to how
    //MegaHAL works to try to choose the most original response from its generations;
    //this is in addition to tyuo's own scoring model, and it's up to the caller
    //to decide which response to display (probably highest scored, tie-broken by
    //highest surprise in most cases)
    //turning either direction on incurs a linear n-gram lookup at the lowest-enabled
    //level, so it may be worth disabling if milliseconds matter
    CalculateSurpriseForward bool
    CalculateSurpriseReverse bool
}
type contextConfig struct {
    Language string //"english", "french"

    Ngrams contextConfigNgrams

    Learning contextConfigLearning

    Production contextConfigProduction
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
func prepareContext(
    contextsPath string,
    contextId string ,
    databaseManager *databaseManager,
    bannedSubstringsGenericByLanguage map[string][]string,
    boringTokensByLanguage map[string]map[string]void,
) (*Context, error) {
    logger.Infof("loading context %s...", contextId)
    
    configFile, err := os.Open(filepath.Join(contextsPath, contextId + ".json"))
    if err != nil {
        logger.Warningf("unable to load context %s: %s", contextId, err)
        return nil, err
    }
    defer configFile.Close()
    
    configJson, err := ioutil.ReadAll(configFile)
    if err != nil {
        return nil, err
    }
    
    var config contextConfig
    if err = json.Unmarshal(configJson, &config); err != nil {
        return nil, err
    }
    
    database, err := databaseManager.Load(contextId)
    if err != nil {
        return nil, err
    }
    
    boringTokens, defined := boringTokensByLanguage[config.Language]
    if !defined {
        return nil, errors.New(fmt.Sprintf("boring tokens not defined for %s", config.Language))
    }
    
    bannedSubstringsGeneric, defined := bannedSubstringsGenericByLanguage[config.Language]
    if !defined {
        return nil, errors.New(fmt.Sprintf("banned tokens not defined for %s", config.Language))
    }
    bannedDictionary, err := prepareBannedDictionary(database, bannedSubstringsGeneric)
    if err != nil {
        return nil, err
    }
    
    dictionary, err := prepareDictionary(database)
    if err != nil {
        return nil, err
    }
    
    return &Context{
        config: config,

        database: database,
        bannedDictionary: bannedDictionary,
        dictionary: dictionary,
        boringTokens: boringTokens,
    }, nil
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

func (c *Context) BanSubstrings(substrings []string) (error) {
    return c.bannedDictionary.ban(stringSliceToSet(substrings))
}
func (c *Context) UnbanSubstrings(substrings []string) (error) {
    return c.bannedDictionary.unban(stringSliceToSet(substrings))
}

func (c *Context) GetProductionTokensInitial() (int) {
    return c.config.Production.TokensInitial
}
func (c *Context) GetMaxParallelOperations() (int) {
    return c.config.Production.MaxParallelOperations
}




func (c *Context) getOldestAllowedTime() (int64) {
    return time.Now().Unix() - c.config.Learning.MaxAge
}

func (c *Context) GetDigrams(
    specs map[DigramSpec]bool,
    forward bool,
) (map[DigramSpec]Digram, error) {
    return c.database.digramsGet(
        specs,
        forward,
        c.getOldestAllowedTime(),
    )
}

func (c *Context) GetTrigrams(
    specs map[TrigramSpec]bool,
    forward bool,
) (map[TrigramSpec]Trigram, error) {
    return c.database.trigramsGet(
        specs,
        forward,
        c.getOldestAllowedTime(),
    )
}
func (c *Context) GetTrigramsOrigin(
    dictionaryIdFirst int,
    count int,
    forward bool,
) ([]Trigram, error) {
    return c.database.trigramsGetOnlyFirst(
        dictionaryIdFirst,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}

func (c *Context) GetQuadgrams(
    specs map[QuadgramSpec]bool,
    forward bool,
) (map[QuadgramSpec]Quadgram, error) {
    return c.database.quadgramsGet(
        specs,
        forward,
        c.getOldestAllowedTime(),
    )
}
func (c *Context) GetQuadgramsOrigin(
    dictionaryIdFirst int,
    count int,
    forward bool,
) ([]Quadgram, error) {
    return c.database.quadgramsGetOnlyFirst(
        dictionaryIdFirst,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}
func (c *Context) GetQuadgramsFromBoundary(
    dictionaryIdSecond int,
    count int,
    forward bool,
) ([]Quintgram, error) {
    return c.database.quadgramsGetFromBoundary(
        dictionaryIdSecond,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}

func (c *Context) GetQuintgrams(
    specs map[QuintgramSpec]bool,
    forward bool,
) (map[QuintgramSpec]Quintgram, error) {
    return c.database.quintgramsGet(
        specs,
        forward,
        c.getOldestAllowedTime(),
    )
}
func (c *Context) GetQuintgramsOrigin(
    dictionaryIdFirst int,
    count int,
    forward bool,
) ([]Quintgram, error) {
    return c.database.quintgramsGetOnlyFirst(
        dictionaryIdFirst,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}
func (c *Context) GetQuintgramsFromBoundary(
    dictionaryIdSecond int,
    count int,
    forward bool,
) ([]Quintgram, error) {
    return c.database.quintgramsGetFromBoundary(
        dictionaryIdSecond,
        count,
        forward,
        c.getOldestAllowedTime(),
    )
}




func (c *Context) IsAllowed(s string) (bool) {
    return !c.bannedDictionary.containsBannedToken(s)
}
func (c *Context) GetIdsBannedStatus(ids []int) (map[int]bool) {
    return c.bannedDictionary.getIdsBannedStatus(intSliceToSet(ids))
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
    depunctuatedTokens = nil //not needed anymore and this function's runtime is far from over

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

        candidates[pt.Base] = false
    }

    ids, err := c.dictionary.getIdsByToken(candidates)
    if err != nil {
        return nil, err
    }
    filteredIds := make([]int, 0, len(ids))
    for _, id := range ids {
        if !c.bannedDictionary.getIdBannedStatus(id) {
            filteredIds = append(filteredIds, id)
        }
    }
    return filteredIds, nil
}





type ContextManager struct {
    contextsPath string
    
    databaseManager *databaseManager

    bannedSubstringsGenericByLanguage map[string][]string
    boringTokensByLanguage map[string]map[string]void

    contexts map[string]*Context

    //used internally to control access to GetContext(), so that
    //resources like the database aren't connected multiple times
    lock sync.Mutex
}
func PrepareContextManager(dataPath string) (*ContextManager, error) {
    bannedSubstringsGenericByLanguage := make(map[string][]string)
    boringTokensByLanguage := make(map[string]map[string]void)
    
    languagesPath := filepath.Join(dataPath, "languages")
    if files, err := ioutil.ReadDir(languagesPath); err != nil {
        return nil, err
    } else {
        for _, file := range files {
            logger.Debugf("evaluating %s...", file.Name())
            if strings.HasSuffix(file.Name(), ".banned") {
                if bannedSubstrings, err := processBannedSubstrings(filepath.Join(languagesPath, file.Name())); err != nil {
                    return nil, err
                } else {
                    bannedSubstringsGenericByLanguage[file.Name()[:len(file.Name()) - 7]] = bannedSubstrings
                }
            } else if strings.HasSuffix(file.Name(), ".boring") {
                if boringTokens, err := processBoringTokens(filepath.Join(languagesPath, file.Name())); err != nil {
                    return nil, err
                } else {
                    boringTokensByLanguage[file.Name()[:len(file.Name()) - 7]] = boringTokens
                }
            }
        }
    }
    
    contextsPath := filepath.Join(dataPath, "contexts")
    return &ContextManager{
        contextsPath: contextsPath,
        
        databaseManager: prepareDatabaseManager(contextsPath),
        
        bannedSubstringsGenericByLanguage: bannedSubstringsGenericByLanguage,
        boringTokensByLanguage: boringTokensByLanguage,
        
        contexts: make(map[string]*Context),
    }, nil
}
func (cm *ContextManager) Close() {
    cm.lock.Lock()
    defer cm.lock.Unlock()
    
    cm.databaseManager.Close()
    cm.contexts = make(map[string]*Context)
}
func (cm *ContextManager) GetContext(contextId string) (*Context, error) {
    cm.lock.Lock()
    defer cm.lock.Unlock()

    if context, defined := cm.contexts[contextId]; defined {
        return context, nil
    }
    
    if context, err := prepareContext(
        cm.contextsPath,
        contextId,
        cm.databaseManager,
        cm.bannedSubstringsGenericByLanguage,
        cm.boringTokensByLanguage,
    ); err == nil {
        cm.contexts[contextId] = context
        return context, nil
    } else {
        return nil, err
    }
}