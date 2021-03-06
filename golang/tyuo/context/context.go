package context

//context-manager holds a bunch of contexts, keyed by ID
//each context has a database file and a language-specifying file as artifacts
//once a context is loaded, a database connection and the language value are
//held in memory

/*
tyuo/
    contexts/
        <id>.sqlite3
        <id>.language
    languages/
        <language>.banned
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
