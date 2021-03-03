use crate::engine::model::banned_dictionary;
use crate::engine::model::dictionary;
use crate::engine::model::model;

use std::error::Error;

use std::collections::{HashMap, HashSet};

use std::io::{Read, Write};

use rusqlite::{Connection, named_params, params};
use serde_json::json;


pub struct Database {
    connection: Connection,
}
impl Database {
    pub fn new(connection:Connection) -> Database {
        return Database{
            connection: connection,
        };
    }


    pub fn banned_load_banned_tokens(&self,
        tokens:Option<HashSet<&str>>,
    ) -> Result<Vec<banned_dictionary::BannedToken>, Box<Error>> {
        let tkns:HashSet<&str>;
        let mut query_string:String = "SELECT
            banned.caseInsensitiveRepresentation,
            dict.id
        FROM
            dictionary_banned AS banned
        LEFT JOIN dictionary AS dict ON
            banned.caseInsensitiveRepresentation = dict.caseInsensitiveRepresentation
        ".to_string();
        if tokens.is_some() {
            tkns = tokens.unwrap();

            let mut array_parms = Vec::with_capacity(tkns.len());
            for idx in 0..tkns.len() {
                array_parms.push(format!("?{}", idx + 1));
            }

            query_string.push_str(format!("WHERE
                banned.caseInsensitiveRepresentation IN ({})
            ", array_parms.join(",")).as_str());
        } else {
            tkns = HashSet::new();
        }
        let mut stmt = self.connection.prepare(query_string.as_str())?;

        let mut results = Vec::new();
        let mut rows = stmt.query(tkns)?;
        while let Some(row) = rows.next()? {
            results.push(banned_dictionary::BannedToken::new(
                row.get(0)?,
                row.get(1)?,
            ));
        }
        return Ok(results);
    }
    pub fn banned_ban_tokens(&mut self,
        tokens:HashSet<&str>,
    ) -> Result<Vec<banned_dictionary::BannedToken>, Box<Error>> {
        let tx = self.connection.transaction()?;
        {
            let mut insert_stmt = tx.prepare("INSERT INTO
                dictionary_banned(caseInsensitiveRepresentation)
            VALUES (?1)
            ON CONFLICT DO NOTHING
            ")?;
            for token in &tokens {
                insert_stmt.execute(&[token])?;
            }
        }
        tx.commit()?;

        return self.banned_load_banned_tokens(Some(tokens));
    }
    pub fn banned_unban_tokens(&self,
        tokens:HashSet<&str>,
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


    pub fn dictionary_enumerate_words_by_substring(&self,
        substrings:HashSet<&str>,
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
    pub fn dictionary_get_words_by_token(&self,
        tokens:HashSet<&str>,
    ) -> Result<Vec<dictionary::DictionaryWord>, Box<Error>> {
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
            results.push(dictionary::DictionaryWord::new(
                row.get(1)?,
                row.get(2)?,
                row.get(0)?,
                map,
            ));
        }
        return Ok(results);
    }
    pub fn dictionary_get_words_by_id(&self,
        ids:HashSet<i32>,
    ) -> Result<Vec<dictionary::DictionaryWord>, Box<Error>> {
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
            results.push(dictionary::DictionaryWord::new(
                row.get(1)?,
                row.get(2)?,
                row.get(0)?,
                map,
            ));
        }
        return Ok(results);
    }
    pub fn dictionary_get_random_words(&self,
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
    pub fn dictionary_set_words(&mut self,
        words:HashSet<dictionary::DictionaryWord>,
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
            for word in words {
                let cfj:Option<Vec<u8>>;
                let capitalised_forms = word.get_capitalised_forms();
                if capitalised_forms.len() > 0 {
                    cfj = Some(serde_json::to_vec(&word.get_capitalised_forms())?);
                } else {
                    cfj = None;
                }
                
                stmt.execute_named(named_params!{
                    ":cir": word.get_case_insensitive_representation(),
                    ":id": word.get_id(),
                    ":cio": word.get_case_insensitive_occurrences(),
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

pub struct DatabaseManager {
    db_dir: std::path::PathBuf,
}
impl DatabaseManager {
    pub fn new(db_dir:&std::path::Path) -> Box<DatabaseManager> {
        return Box::new(DatabaseManager{
            db_dir: db_dir.to_owned(),
        });
    }

    fn resolve_path(&self, id:&str) -> std::path::PathBuf {
        let mut db_path = self.db_dir.join(id);
        db_path.set_extension("sqlite3");
        return db_path;
    }
    pub fn exists(&self, id:&str) -> bool {
        return self.resolve_path(id).is_file();
    }
    pub fn load(&self, id:&str) -> Result<Database, Box<Error>> {
        let path = self.resolve_path(id);
        info!("loading database {}...", id);
        let connection = Connection::open(path)?;

        debug!("preparing database structures...");
        connection.execute("CREATE TABLE IF NOT EXISTS dictionary (
            caseInsensitiveRepresentation TEXT NOT NULL UNIQUE,
            id INTEGER NOT NULL PRIMARY KEY,
            caseInsensitiveOccurrences INTEGER NOT NULL,
            capitalisedFormsJSON BLOB
        )", params![])?;
        connection.execute("CREATE TABLE IF NOT EXISTS dictionary_banned (
            caseInsensitiveRepresentation TEXT NOT NULL PRIMARY KEY
        )", params![])?;
        //for the next two, the JSON structure will never be empty, since there
        //has to be at least one transition for a write to occur
        connection.execute("CREATE TABLE IF NOT EXISTS statistics_forward (
            dictionaryId INTEGER NOT NULL PRIMARY KEY,
            childrenJSONZLIB BLOB NOT NULL,

            FOREIGN KEY(dictionaryId) REFERENCES dictionary(id) ON DELETE CASCADE
        )", params![])?;
        connection.execute("CREATE TABLE IF NOT EXISTS statistics_reverse (
            dictionaryId INTEGER NOT NULL PRIMARY KEY,
            childrenJSONZLIB BLOB NOT NULL,

            FOREIGN KEY(dictionaryId) REFERENCES dictionary(id) ON DELETE CASCADE
        )", params![])?;

        debug!("preparing database pragma...");
        //ensure foreign keys are checked
        connection.execute("PRAGMA foreign_keys=on", params![])?;
        //since only lower-case matches occur, let comparisons be optimal
        connection.execute("PRAGMA case_sensitive_like=true", params![])?;

        return Ok(Database::new(connection));
    }
    pub fn drop(&self, id:&str) -> Result<(), String> {
        if std::fs::remove_file(self.resolve_path(id)).is_ok() {
            return Ok(());
        }
        return Err(format!("unable to remove {}", id));
    }
}
