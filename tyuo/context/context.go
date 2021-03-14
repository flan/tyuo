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


type contextConfigNgrams struct {
    Digrams bool
    Trigrams bool
    Quadgrams bool
    Quintgrams bool
}
type contextConfigLearning struct {
    MinTokenCount int
    
    MaxTokenLength int

    MaxAge int64

    RescaleThreshold int
    RescaleDecimator int
}
type contextConfigProduction struct {
    MaxParallelSearches int
    
    TokensInitial int
    SearchBranchesInitial int
    SearchBranchesFromBoundaryInitial int
    SearchBranchesChildren int

    MinLength int
    MaxLength int
    StopProbability float32

    TargetMinLength int
    TargetMaxLength int
    TargetStopProbability float32

    BaseRepresentationThreshold float32
    
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
func (c *Context) GetProductionSearchBranchesInitial() (int) {
    return c.config.Production.SearchBranchesInitial
}
func (c *Context) GetProductionSearchBranchesFromBoundaryInitial() (int) {
    return c.config.Production.SearchBranchesFromBoundaryInitial
}
func (c *Context) GetProductionSearchBranchesChildren() (int) {
    return c.config.Production.SearchBranchesChildren
}

func (c *Context) GetProductionMaxParallelSearches() (int) {
    return c.config.Production.MaxParallelSearches
}

func (c *Context) GetProductionMinLength() (int) {
    return c.config.Production.MinLength
}
func (c *Context) GetProductionMaxLength() (int) {
    return c.config.Production.MaxLength
}
func (c *Context) GetProductionStopProbability() (float32) {
    return c.config.Production.StopProbability
}
func (c *Context) GetProductionTargetMinLength() (int) {
    return c.config.Production.TargetMinLength
}
func (c *Context) GetProductionTargetMaxLength() (int) {
    return c.config.Production.TargetMaxLength
}
func (c *Context) GetProductionTargetStopProbability() (float32) {
    return c.config.Production.TargetStopProbability
}

func (c *Context) GetProductionBaseRepresentationThreshold() (float32) {
    return c.config.Production.BaseRepresentationThreshold
}

func (c *Context) GetProductionCalculateSurpriseForward() (bool) {
    return c.config.Production.CalculateSurpriseForward
}
func (c *Context) GetProductionCalculateSurpriseReverse() (bool) {
    return c.config.Production.CalculateSurpriseReverse
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
) ([]Quadgram, error) {
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
func (c *Context) AreIdsAllowed(ids []int) (bool) {
    for _, banned := range c.GetIdsBannedStatus(ids) {
        if banned {
            return false
        }
    }
    return true
}

func (c *Context) LearnInput(tokens []ParsedToken) (error) {
    if len(tokens) < c.config.Learning.MinTokenCount {
        return nil
    }

    rescaleThreshold := c.config.Learning.RescaleThreshold
    rescaleDecimator := c.config.Learning.RescaleDecimator

    //strip punctuation and symbols so they don't get learned redundantly
    deReservedTokens := make([]ParsedToken, 0, len(tokens))
    for _, pt := range tokens {
        if _, defined := PunctuationIdsByToken[pt.Base]; !defined {
            if _, defined = SymbolsIdsByToken[pt.Base]; !defined {
                deReservedTokens = append(deReservedTokens, pt)
            }
        }
    }

    //first, update the dictionary to make sure all tokens have an ID
    dictionaryTokens, err := c.dictionary.learnTokens(
        deReservedTokens,
        rescaleThreshold,
        rescaleDecimator,
    )
    if err != nil {
        return err
    }
    deReservedTokens = nil //not needed anymore and this function's runtime is far from over

    tokensMap := make(map[string]int, len(dictionaryTokens) + len(PunctuationIdsByToken))
    for _, dt := range dictionaryTokens {
        tokensMap[dt.baseRepresentation] = dt.id
    }
    //put punctuation mappings in
    for token, id := range PunctuationIdsByToken {
        tokensMap[token] = id
    }
    //put symbols mappings in
    for token, id := range SymbolsIdsByToken {
        tokensMap[token] = id
    }
    
    //correctness check, to avoid a class of error where a stripped token-subset isn't restored
    //if this occurs, the database could end up corrupted, which is very not-good
    for _, dt := range dictionaryTokens {
        if _, defined := tokensMap[dt.baseRepresentation]; !defined {
            return errors.New(fmt.Sprintf("unable to find a dictionary binding for %s", dt.baseRepresentation))
        }
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

func (c *Context) GetDictionaryTokensById(ids map[int]bool) (map[int]DictionaryToken, error) {
    return c.dictionary.getSliceById(ids)
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
