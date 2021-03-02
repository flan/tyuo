use crate::engine::model::banned_dictionary;
use crate::engine::model::dictionary;

use std::error::Error;

use std::collections::{HashMap, HashSet};

use std::io::{Read, Write};

use rusqlite::{Connection, params};
use serde_json::json;

fn to_json_zlib(serialisable_data:serde_json::Value) -> Vec<u8> {
    let mut encoder = flate2::write::ZlibEncoder::new(Vec::new(), flate2::Compression::default());
    encoder.write_all(&serde_json::to_vec(&serialisable_data).unwrap());
    return encoder.finish().unwrap();
}
fn from_json_zlib(blob:Vec<u8>) -> Option<serde_json::Value> {
    let mut decoder = flate2::read::ZlibDecoder::new(blob.as_slice());
    let mut decoded = String::new();
    if decoder.read_to_string(&mut decoded).is_err() {
        return None;
    }
    
    let deserialised_json = serde_json::from_str(decoded.as_str());
    if deserialised_json.is_ok() {
        return Some(deserialised_json.unwrap());
    }
    return None;
}


pub struct Database {
    connection: Connection,
}
impl Database {
    pub fn prepare(connection:Connection) -> Database {
        return Database{
            connection: connection,
        };
    }
    
    
    pub fn banned_load_banned_tokens(&mut self, tokens:Option<Vec<&str>>) -> Result<Vec<banned_dictionary::BannedWord>, Box<Error>> {
        let tkns:Vec<&str>;
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
            tkns = Vec::new();
        }
        let mut stmt = self.connection.prepare(query_string.as_str())?;
        
        let mut results = Vec::new();
        for _ in stmt.query_map(tkns, |row| {
            results.push(banned_dictionary::BannedWord::prepare(
                row.get(0)?,
                row.get(1)?,
            ));
            return Ok(());
        })?{};
        
        return Ok(results);
    }
    pub fn banned_ban_tokens(&mut self, tokens:Vec<&str>) -> Result<Vec<banned_dictionary::BannedWord>, Box<Error>> {
        let tx = self.connection.transaction()?;
        {
            let mut insert_stmt = tx.prepare("INSERT INTO
                dictionary_banned(caseInsensitiveRepresentation)
            VALUES (?1)
            ON CONFLICT DO NOTHING
            ")?;
            for token in tokens.to_vec() { //copy it, since there's another use later
                insert_stmt.execute(&[token])?;
            }
        }
        tx.commit()?;
        
        return self.banned_load_banned_tokens(Some(tokens));
    }
    pub fn banned_unban_tokens(&self, tokens:Vec<&str>) -> Result<(), Box<Error>> {
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
    
    
    pub fn dictionary_enumerate_words_by_substring(&self, substrings:Vec<&str>) -> Result<HashMap<String, i32>, Box<Error>> {
        
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
            for _ in stmt.query_map(&[format!("%{}%", substring)], |row| {
                results.insert(row.get(0)?, row.get(1)?);
                return Ok(());
            })?{};
        }
        
        return Ok(results);
    }
    pub fn dictionary_get_words_by_token(&self) {
        //exact-match
        //creates if not defined, including upsert
    }
    pub fn dictionary_get_words_by_id(&self, ids:HashSet<i32>) -> Result<Vec<dictionary::DictionaryWord>, Box<Error>> {
        let mut array_parms = Vec::with_capacity(ids.len());
        let ids_len = ids.len();
        for idx in 0..ids_len {
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
        
        let mut results = Vec::with_capacity(ids_len);
        for _ in stmt.query_map(ids, |row| {
            let raw_json:Vec<u8> = row.get(3)?;
            let raw_json_str = std::str::from_utf8(&raw_json)?;
            let deserialised_json = serde_json::from_str(raw_json_str).unwrap();
            
            results.push(dictionary::DictionaryWord::prepare(
                row.get(1)?,
                row.get(2)?,
                row.get(0)?,
                deserialised_json,
            ));
            return Ok(());
        })?{};
        
        if results.len() != ids_len {
            return Err(string_error::into_err(format!(
                "{} requested IDs were not found",
                ids_len - results.len(),
            )));
        }
        
        return Ok(results);
    }
    pub fn dictionary_set_words(&self) {
        //the level consuming this should delete case-sensitive variations of banned words
        //to save a bit of space
    }
    pub fn dictionary_get_latest_identifier(&self) -> Result<i32, Box<Error>>{
        return Ok(0);
    }
    
    pub fn model_get_transitions(&self, direction:&str) {
        
    }
    pub fn model_set_transitions(&self, direction:&str) {
        //if the node has no transitions, delete it; this efficiently reinforces the banned case
    }
}

pub struct DatabaseManager {
    db_dir: std::path::PathBuf,
}
impl DatabaseManager {
    pub fn prepare(db_dir:&std::path::Path) -> Box<DatabaseManager> {
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
        connection.execute("CREATE TABLE IF NOT EXISTS statistics_forward (
            dictionaryId INTEGER NOT NULL PRIMARY KEY,
            childrenJSONZLIB BLOB,
            
            FOREIGN KEY(dictionaryId) REFERENCES dictionary(id) ON DELETE CASCADE
        )", params![])?;
        connection.execute("CREATE TABLE IF NOT EXISTS statistics_reverse (
            dictionaryId INTEGER NOT NULL PRIMARY KEY,
            childrenJSONZLIB BLOB,
            
            FOREIGN KEY(dictionaryId) REFERENCES dictionary(id) ON DELETE CASCADE
        )", params![])?;
        
        debug!("preparing database pragma...");
        //ensure foreign keys are checked
        connection.execute("PRAGMA foreign_keys=on", params![])?;
        //since only lower-case matches occur, let comparisons be optimal
        connection.execute("PRAGMA case_sensitive_like=true", params![])?;
        
        return Ok(Database::prepare(connection));
        /*
        let e = connection.err();
        if e.is_some() {
            error!("unable to load database {}: {}", id, e.unwrap());
        } else {
            error!("unable to load database {}", id);
        }
        return Err("unable to open database");
        *///TODO: reuse this error-presentation logic higher up
    }
    pub fn drop(&self, id:&str) -> Result<(), String> {
        if std::fs::remove_file(self.resolve_path(id)).is_ok() {
            return Ok(());
        }
        return Err(format!("unable to remove {}", id));
    }
}
