package context
import (
    "database/sql"
    "path/filepath"
    "fmt"
    "os"
    "strings"

    _ "github.com/mattn/go-sqlite3"
)

func prepareSqliteArrayParams(start int32, count int32) (string) {
    arrayParams := make([]string, count)
    for i := 0; i < count; i++ {
        arrayParams[i] = fmt.Sprintf("?%d", start + i)
    }
    return strings.Join(arrayParams, ",")
}

func stringSliceToInterfaceSlice(s []string) ([]interface{}) {
    output := make([]interface{}, len(s))
    for i, v := range(s) {
        output[i] = v
    }
    return output
}

type Database struct {
    connection *sql.DB
}
func prepareDatabase(
    dbPath string,
) (*Database, error) {
    connection, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    logger.Debugf("preparing database structures...");
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS dictionary (
        caseInsensitiveRepresentation TEXT NOT NULL UNIQUE,
        id INTEGER NOT NULL PRIMARY KEY,
        caseInsensitiveOccurrences INTEGER NOT NULL,
        capitalisedFormsJSON BLOB
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS dictionary_banned (
        caseInsensitiveRepresentation TEXT NOT NULL PRIMARY KEY
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    //for n-grams, the JSON structure will never be empty, since there
    //has to be at least one transition for a write to occur
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS terminals (
        dictionaryId INTEGER NOT NULL PRIMARY KEY,
        lastObservedForward INTEGER, --UNIX timestamp
        lastObservedReverse INTEGER, --UNIX timestamp

        FOREIGN KEY(dictionaryId)
        REFERENCES dictionary(id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS trigrams_forward (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        childrenJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond)
        REFERENCES dictionary(id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS trigrams_reverse (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        childrenJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond)
        REFERENCES dictionary(id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS quadgrams_forward (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        dictionaryIdThird INTEGER NOT NULL,
        childrenJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird)
        REFERENCES dictionary(id, id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS quadgrams_reverse (
        dictionaryIdFirst INTEGER NOT NULL,
        dictionaryIdSecond INTEGER NOT NULL,
        dictionaryIdThird INTEGER NOT NULL,
        childrenJSONZLIB BLOB NOT NULL,

        PRIMARY KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird),
        FOREIGN KEY(dictionaryIdFirst, dictionaryIdSecond, dictionaryIdThird)
        REFERENCES dictionary(id, id, id)
        ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    logger.Debugf("preparing database pragma...");
    //while foreign keys are declared in the structure, because tokens are never
    //removed from the database, their enforcement is unnecessary
    if _, err = connection.Exec("PRAGMA foreign_keys = OFF"); err != nil {
        connection.Close()
        return nil, err
    }
    //since only lower-case matches occur, make comparisons more efficient
    if _, err = connection.Exec("PRAGMA case_sensitive_like = TRUE"); err != nil {
        connection.Close()
        return nil, err
    }

    return &Database{
        connection: connection,
    }, nil
}
func (db *Database) Close() (error) {
    return db.connection.Close()
}
func (db *Database) bannedLoadBannedTokens(
    tokenSubset []string,
) ([]bannedToken, error) {
    query := `SELECT
        banned.caseInsensitiveRepresentation,
        dict.id
    FROM
        dictionary_banned AS banned
    LEFT JOIN dictionary AS dict ON
        banned.caseInsensitiveRepresentation = dict.caseInsensitiveRepresentation
    `
    
    if len(tokenSubset) > 0 {
        query += fmt.Sprintf(
            "WHERE banned.caseInsensitiveRepresentation IN (%s)",
            prepareSqliteArrayParams(1, len(tokenSubset)),
        )
    }
    if rows, err := db.connection.Query(
        query,
        stringSliceToInterfaceSlice(tokenSubset)...,
    ); err == nil {
        defer rows.Close()
        
        output := make([]bannedToken, 0)
        for rows.Next() {
            var cir string
            var did int32
            if err:= rows.Scan(&cir, &did); err != nil {
                output = append(output, bannedToken{
                    caseInsensitiveRepresentation: cir,
                    dictionaryId: did,
                })
            } else {
                return nil, err
            }
        }
        return output, nil
    } else {
        return nil, err
    }
}
func (db *Database) bannedBanTokens(tokens []string) ([]bannedToken, error) {
    tx, err := db.connection.Begin()
    if err != nil {
        return nil, err
    }
    
    const query = `INSERT INTO
        dictionary_banned(caseInsensitiveRepresentation)
    VALUES (?1)
    ON CONFLICT DO NOTHING
    `
    if stmt, err := tx.Prepare(query); err == nil {
        for _, token := range tokens {
            if _, err = stmt.Exec(token); err != nil {
                break
            }
        }
        stmt.Close()
    }
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    if err = tx.Commit(); err != nil {
        return nil, err
    }
    return db.bannedLoadBannedTokens(tokens);
}
func (db *Database) bannedUnbanTokens(tokens []string) (error) {
    query := fmt.Sprintf(`DELETE
    FROM
        dictionary_banned
    WHERE caseInsensitiveRepresentation IN (%s)
    `, prepareSqliteArrayParams(1, len(tokens)))
    
    _, err := db.connection.Exec(query, stringSliceToInterfaceSlice(tokens)...)
    return err
}





/*
impl Database {

    pub fn dictionary_enumerate_tokens_by_substring(&self,
        substrings:&HashSet<&str>,
    ) -> Result<HashMap<String, i32>, Box<Error>> {
        let mut stmt = self.connection.prepare("SELECT
            caseInsensitiveRepresentation,
            id
        FROM
            dictionary
        WHERE
            caseInsensitiveRepresentation LIKE ?1
        ")?;

        let mut results = HashMap::new();
        for substring in substrings {
            let mut rows = stmt.query(&[format!("%{}%", substring)])?;
            while let Some(row) = rows.next()? {
                results.insert(row.get(0)?, row.get(1)?);
            }
        }
        return Ok(results);
    }
    pub fn dictionary_get_tokens_by_token(&self,
        tokens:&HashSet<&str>,
    ) -> Result<Vec<dictionary::DictionaryToken>, Box<Error>> {
        let mut array_parms = Vec::with_capacity(tokens.len());
        for idx in 0..tokens.len() {
            array_parms.push(format!("?{}", idx + 1));
        }

        let mut stmt = self.connection.prepare(format!("SELECT
            caseInsensitiveRepresentation,
            id,
            caseInsensitiveOccurrences,
            capitalisedFormsJSON
        FROM
            dictionary
        WHERE
            caseInsensitiveRepresentation IN ({})
        ", array_parms.join(",")).as_str())?;

        let mut results = Vec::with_capacity(tokens.len());
        let mut rows = stmt.query(tokens)?;
        while let Some(row) = rows.next()? {
            let map:HashMap<String, i32>;
            let raw_json:Option<Vec<u8>> = row.get(3)?;
            match raw_json {
                Some(data) => map = serde_json::from_str(
                    std::str::from_utf8(&data)?,
                )?,
                None => map = HashMap::new(),
            }
            results.push(dictionary::DictionaryToken::new(
                row.get(1)?,
                row.get(2)?,
                row.get(0)?,
                map,
            ));
        }
        return Ok(results);
    }
    pub fn dictionary_get_tokens_by_id(&self,
        ids:&HashSet<i32>,
    ) -> Result<Vec<dictionary::DictionaryToken>, Box<Error>> {
        let mut array_parms = Vec::with_capacity(ids.len());
        for idx in 0..ids.len() {
            array_parms.push(format!("?{}", idx + 1));
        }

        let mut stmt = self.connection.prepare(format!("SELECT
            caseInsensitiveRepresentation,
            id,
            caseInsensitiveOccurrences,
            capitalisedFormsJSON
        FROM
            dictionary
        WHERE
            id IN ({})
        ", array_parms.join(",")).as_str())?;

        let mut results = Vec::with_capacity(ids.len());
        let mut rows = stmt.query(ids)?;
        while let Some(row) = rows.next()? {
            let map:HashMap<String, i32>;
            let raw_json:Option<Vec<u8>> = row.get(3)?;
            match raw_json {
                Some(data) => map = serde_json::from_str(
                    std::str::from_utf8(&data)?,
                )?,
                None => map = HashMap::new(),
            }
            results.push(dictionary::DictionaryToken::new(
                row.get(1)?,
                row.get(2)?,
                row.get(0)?,
                map,
            ));
        }
        return Ok(results);
    }
    pub fn dictionary_set_tokens(&mut self,
        tokens:HashSet<dictionary::DictionaryToken>,
    ) -> Result<(), Box<Error>> {
        let tx = self.connection.transaction()?;
        {
            let mut stmt = tx.prepare("INSERT INTO dictionary(
                caseInsensitiveRepresentation,
                id,
                caseInsensitiveOccurrences,
                capitalisedFormsJSON
            ) VALUES (:cir, :id, :cio, :cfj)
            ON CONFLICT(id) DO UPDATE SET
                caseSensitiveOccurrences = :cio,
                capitalisedFormsJSON = :cfj
            ")?;
            for token in tokens {
                let cfj:Option<Vec<u8>>;
                let capitalised_forms = token.get_capitalised_forms();
                if capitalised_forms.len() > 0 {
                    cfj = Some(serde_json::to_vec(&token.get_capitalised_forms())?);
                } else {
                    cfj = None;
                }
                
                stmt.execute_named(named_params!{
                    ":cir": token.get_case_insensitive_representation(),
                    ":id": token.get_id(),
                    ":cio": token.get_case_insensitive_occurrences(),
                    ":cfj": cfj,
                })?;
            }
        }
        tx.commit()?;
        return Ok(());
    }
    pub fn dictionary_get_next_identifier(&self) -> Result<i32, Box<Error>>{
        let mut next_identifier:i32 = -2147483648; //lowest allowable identifier
        self.connection.query_row("SELECT MAX(id) FROM dictionary", params![], |row| {
            next_identifier = row.get(0)?;
            return Ok(());
        })?;
        return Ok(next_identifier);
    }


    pub fn model_get_transitions(&self,
        direction:&str,
        ids:HashSet<&i32>,
    ) -> Result<fnv::FnvHashMap<i32, fnv::FnvHashMap<i32, model::Transition>>, Box<Error>> {
        let mut array_parms = Vec::with_capacity(ids.len());
        for idx in 0..ids.len() {
            array_parms.push(format!("?{}", idx + 1));
        }

        let mut stmt = self.connection.prepare(format!("SELECT
            dictionaryId,
            childrenJSONZLIB
        FROM
            statistics_{}
        WHERE
            dictionaryId IN ({})
        ", direction, array_parms.join(",")).as_str())?;

        let mut results = fnv::FnvHashMap::default();
        let mut rows = stmt.query(ids)?;
        while let Some(row) = rows.next()? {
            let id:i32 = row.get(0)?;
            
            let blob:Vec<u8> = row.get(1)?;
            let mut decoder = flate2::read::ZlibDecoder::new(blob.as_slice());
            let mut decoded = String::new();
            decoder.read_to_string(&mut decoded)?;
            
            let mut transitions = fnv::FnvHashMap::default();
            let vec:Vec<serde_json::Value> = serde_json::from_str(decoded.as_str())?;
            for transition in vec {
                let dictionary_id:i32 = transition[0].as_i64().unwrap() as i32;
                let occurrences:i32 = transition[1].as_i64().unwrap() as i32;
                let last_observed:i64 = transition[1].as_i64().unwrap();
                
                transitions.insert(dictionary_id, model::Transition::new(
                    occurrences,
                    last_observed,
                ));
            }
            
            results.insert(id, transitions);
        }
        return Ok(results);
    }
    pub fn model_set_transitions(&mut self,
        direction:&str,
        nodes:Vec<(&i32, HashMap<&i32, model::Transition>)>,
    ) -> Result<(), Box<Error>> {
        let tx = self.connection.transaction()?;
        {
            let mut stmt = tx.prepare(format!("INSERT INTO statistics_{}(
                dictionaryId,
                childrenJSONZLIB
            ) VALUES (:id, :cjz)
            ON CONFLICT(dictionaryId) DO UPDATE SET
                childrenJSONZLIB = :cjz
            ", direction).as_str())?;
            for node in nodes {
                let mut transitions:Vec<(i32, i32, i64)> = Vec::with_capacity(node.1.len());
                
                for (transition_dictionary_id, transition) in node.1 {
                    transitions.push((
                        transition_dictionary_id.to_owned(),
                        transition.get_occurrences().to_owned(),
                        transition.get_last_observed().to_owned(),
                    ));
                }
                let json_data = serde_json::json!(transitions);
                let serialisable_data = &serde_json::to_vec(&json_data)?;
                
                let mut encoder = flate2::write::ZlibEncoder::new(Vec::new(), flate2::Compression::default());
                encoder.write_all(serialisable_data);
                
                stmt.execute_named(named_params!{
                    ":id": node.0,
                    ":cjz": encoder.finish()?,
                })?;
            }
        }
        tx.commit()?;
        return Ok(());
    }
}
*/

//NOTE: terminal, trigam, and quadgram logic is needed instead of the
//current digram-only
//when looking up trigrams, input is a tuple; for quadgrams, it's a triple
//trigrams and quadgrams also need "Only1" variants for their
//lookups, which allows wildcard logic when selecting from the database
//(only match on the first column), used to start a search from the
//initial keyword state
//all of the n-gram lookups should use a combination of
//ORDER BY RANDOM() and LIMIT X

//terminal lookup's response is a pair of bools, indicating whether it is
//recognised as a forward or reverse terminal

//there also needs to be a function to select a few reverse-terminals for
//use as a starting point for beginning a random walk as a fallback
//for production flows.

//multiple requests can be made at once (using a prepared statement approach);
//the values returned will be in the same order as they were received













type DatabaseManager struct {
    dbDir string
    
    databases map[string]*Database
}
func prepareDatabaseManager(dbDir string) (*DatabaseManager) {
    return &DatabaseManager{
        dbDir: dbDir,
        
        databases: make(map[string]*Database),
    }
}
func (dbm *DatabaseManager) Close() {
    logger.Debugf("closing databases...")
    for _, database := range dbm.databases {
        database.Close();
    }
    dbm.databases = make(map[string]*Database)
}
func (dbm *DatabaseManager) idToPath(contextId string) (string) {
    return filepath.Join(dbm.dbDir, contextId + ".sqlite3")
}
func (dbm *DatabaseManager) Create(contextId string) (error) {
    logger.Infof("creating database {}...", contextId)
    if db, err := prepareDatabase(dbm.idToPath(contextId)); err == nil {
        return db.Close()
    } else {
        return err
    }
}
func (dbm *DatabaseManager) Drop(contextId string) (error) {
    logger.Infof("dropping database {}...", contextId)
    if database, defined := dbm.databases[contextId]; defined {
        if err := database.Close(); err != nil {
            //it'll be referenced anyway, so this isn't critical
            logger.Warningf("unable to close database %s: %s", contextId, err)
        }
        delete(dbm.databases, contextId)
    }
    return os.Remove(dbm.idToPath(contextId))
}
func (dbm *DatabaseManager) Load(contextId string) (*Database, error) {
    logger.Infof("loading database {}...", contextId)
    dbPath := dbm.idToPath(contextId)
    if _, err := os.Stat(dbPath); err != nil {
        return nil, err
    }
    return prepareDatabase(dbPath)
}
