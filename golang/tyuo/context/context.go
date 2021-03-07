package context

//context-manager holds a bunch of contexts, keyed by ID
//each context has a database file and a language-specifying file as artifacts
//once a context is loaded, a database connection and the language value are
//held in memory

/*
tyuo/
    contexts/
        <id>.sqlite3
        <id>.json {
            "language": "en",
            "n-grams": {
                "digrams": false,
                "trigrams": true,
                "quadgrams": true,
                "quintgrams": true,
            },
        }
    languages/
        <language>.banned
        <language>.nonkey
*/

func Test(contextDir string) () {
    dbm := prepareDatabaseManager(contextDir)
    if err := dbm.Create("hi"); err != nil {
        logger.Errorf("unable to create database: %s", err)
    }
    if err := dbm.Drop("hi"); err != nil {
        logger.Errorf("unable to drop database: %s", err)
    }
}


type Context struct {
    language string
    
    digrams bool
    trigrams bool
    quadgrams bool
    quintgrams bool
}
func (c *Context) Language() (string) {
    return c.language
}

//digrams are the simplest and fastest transition model; using them will
//produce behaviour that is often novel, sometimes insightful, frequently
//deranged, particularly as learning progresses.
//It's pretty random and will only resemble native speech by chance.
func (c *Context) DigramsEnabled() (bool) {
    return c.digrams
}
//trigrams are a fairly middle-ground option, producing relevant observations
//with some regularity and having about as much sentence-structure correctness
//as a machine-translation from a language with no common ancestry.
func (c *Context) TrigramsEnabled() (bool) {
    return c.trigrams
}
//quadgrams are a reasonably stable choice for production of "how do you do,
//fellow humans" responses, being sell-formed for the most part, but closely
//reflecting observed input: a lot of data will need to be learned before novel
//structures will be produced with any regularity and search-spaces will
//sometimes be exhausted while finding viable paths.
func (c *Context) QuadgramsEnabled() (bool) {
    return c.quadgrams
}
//quintgrams (and anything above them) will rarely deviate from mimicing
//what was learned; occasional novel productions are possible, but it will
//not be uncommon to see near-verbatim recreations of input data.
func (c *Context) QuintgramsEnabled() (bool) {
    return c.quintgrams
}
//there might not be a need to give the database an explicit create step
//if there's a JSON file, that's probably good enough to proceed
//there also needs to be a Lock and Unlock function; these are what
//allow the TCP service to not care how many requests it serves
//and to control simultaneous access to the database

type ContextManager struct {
    contexts map[string]Context
}
//functions to create, get, and drop contexts
//create needs to ensure the given language is recognised

