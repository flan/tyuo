package context;
import (
    "database/sql"
    "filepath"
    "fmt"
    "os"
    "os/user"

    _ "github.com/mattn/go-sqlite3"
)


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
    //for the next two, the JSON structure will never be empty, since there
    //has to be at least one transition for a write to occur
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS statistics_forward (
        dictionaryId INTEGER NOT NULL PRIMARY KEY,
        childrenJSONZLIB BLOB NOT NULL,

        FOREIGN KEY(dictionaryId) REFERENCES dictionary(id) ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    if _, err = connection.Exec(`CREATE TABLE IF NOT EXISTS statistics_reverse (
        dictionaryId INTEGER NOT NULL PRIMARY KEY,
        childrenJSONZLIB BLOB NOT NULL,

        FOREIGN KEY(dictionaryId) REFERENCES dictionary(id) ON DELETE CASCADE
    )`); err != nil {
        connection.Close()
        return nil, err
    }
    
    logger.Debugf("preparing database pragma...");
    //ensure foreign keys are checked
    if _, err = connection.Exec("PRAGMA foreign_keys=on"); err != nil {
        connection.Close()
        return nil, err
    }
    //since only lower-case matches occur, let comparisons be optimal
    if _, err = connection.Exec("PRAGMA case_sensitive_like=true"); err != nil {
        connection.Close()
        return nil, err
    }

    return &Database{
        connection: connection,
    }
}
func (db *Database) Close() {
    db.connection.Close()
}
func (db *Database) bannedLoadBannedTokens(
    tokenSubset map[string]bool,
) ([]BannedToken, error) {
    query := `SELECT
        banned.caseInsensitiveRepresentation,
        dict.id
    FROM
        dictionary_banned AS banned
    LEFT JOIN dictionary AS dict ON
        banned.caseInsensitiveRepresentation = dict.caseInsensitiveRepresentation
    `
    
    args := make([]string, len(tokenSubset))
    if len(args) > 0 {
        array_params := make([]string, len(args))
        for i, arg := args {
            args[i] = arg
            array_params[i] = fmt.Sprintf("?{}", i + 1)
        }
        query += fmt.Sprintf(
            "WHERE banned.caseInsensitiveRepresentation IN ({})",
            strings.Join(array_parms, ","),
        )
    }
    if rows, err := db.connection.Query(query, args...); err == nil {
        defer rows.Close()
        
        output := make([]BannedToken, 0)
        for rows.Next() {
            var cir string
            var did int
            if err:= rows.Scan(&cir, &did); err != nil {
                output = append(output, BannedToken{
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









impl Database {

    pub fn banned_ban_tokens(&mut self,
        tokens:&HashSet<&str>,
    ) -> Result<Vec<banned_dictionary::BannedToken>, Box<Error>> {
        let tx = self.connection.transaction()?;
        {
            let mut insert_stmt = tx.prepare("INSERT INTO
                dictionary_banned(caseInsensitiveRepresentation)
            VALUES (?1)
            ON CONFLICT DO NOTHING
            ")?;
            for token in tokens {
                insert_stmt.execute(&[token])?;
            }
        }
        tx.commit()?;

        return self.banned_load_banned_tokens(Some(tokens));
    }
    pub fn banned_unban_tokens(&self,
        tokens:&HashSet<&str>,
    ) -> Result<(), Box<Error>> {
        let mut array_parms = Vec::with_capacity(tokens.len());
        for idx in 0..tokens.len() {
            array_parms.push(format!("?{}", idx + 1));
        }

        self.connection.execute(format!("DELETE
        FROM
            dictionary_banned
        WHERE
            caseInsensitiveRepresentation IN ({})
        ", array_parms.join(",")).as_str(), tokens)?;
        return Ok(());
    }


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
    pub fn dictionary_get_random_tokens(&self,
        count:u8,
    ) -> Result<Vec<(String, i32)>, Box<Error>> {
        let mut stmt = self.connection.prepare(format!("SELECT
            caseInsensitiveRepresentation,
            id
        FROM
            dictionary
        ORDER BY RANDOM()
        LIMIT {}
        ", count).as_str())?;

        let mut results = Vec::with_capacity(count.into());
        let mut rows = stmt.query(params![])?;
        while let Some(row) = rows.next()? {
            results.push((row.get(0)?, row.get(1)?));
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
















type DatabaseManager struct {
    dbDir: string,
    
    databases: map[string]*Database,
}
func PrepareDatabaseManager(dbDir string) (*DatabaseManager) {
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
    return filepath.Join(dbm.dbDir, contextId + '.sqlite3')
}
func (dbm *DatabaseManager) Create(contextId string) {
    logger.Infof("creating database {}...", contextId)
    prepareDatabase(dbm.idToPath(contextId))
}
func (dbm *DatabaseManager) Drop(contextId string) {
    logger.Infof("dropping database {}...", contextId)
    if database, defined := dbm.databases[contextId]; defined {
        database.Close()
        delete(dbm.databases, contextId)
    }
    if err := os.Remove(dbm.idToPath(contextId)); err != nil {
        logger.Warnf("unable to unlink database {}: {}", contextId, e)
    }
}
func (dbm *DatabaseManager) Load(contextId string) (*Database, error) {
    logger.Infof("loading database {}...", contextId)
    dbPath := dbm.idToPath(contextId)
    if _, err := os.Stat(dbPath); err {
        return nil, err
    }
    return prepareDatabase(dbPath)
}
