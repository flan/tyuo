use crate::engine::model::banned_dictionary;

use std::io::{Read, Write};
//use itertools::Itertools;

use rusqlite::{Connection, Error, params};
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
    
    
    pub fn banned_load_banned_tokens(&mut self, tokens:Option<Vec<&str>>) -> Result<Vec<banned_dictionary::BannedWord>, Error> {
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
            
            let mut array_parms = Vec::new();
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
        let iter = stmt.query_map(tkns, |row| {
            return Ok(banned_dictionary::BannedWord::prepare(
                row.get(0)?,
                row.get(1)?,
            ));
        })?;
        
        let mut output = Vec::new();
        for banned_word in iter {
            output.push(banned_word.unwrap());
        }
        return Ok(output);
    }
    pub fn banned_ban_tokens(&mut self, tokens:Vec<&str>) -> Result<Vec<banned_dictionary::BannedWord>, Error> {
        let tx = self.connection.transaction()?;
        {
            let mut insert_stmt = tx.prepare("INSERT INTO
                dictionary_banned(caseInsensitiveRepresentation)
            VALUES (?1)
            ON CONFLICT DO NOTHING
            ")?;
            for token in tokens.to_vec() { //copy it, since there's another use later
                insert_stmt.execute(&[token]);
            }
        }
        tx.commit()?;
        
        return self.banned_load_banned_tokens(Some(tokens));
    }
    pub fn banned_unban_tokens(&self, tokens:Vec<&str>) -> Result<(), Error> {
        let mut array_parms = Vec::new();
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
    
    
    pub fn dictionary_enumerate_words_by_token(&self) {
        //substring-match
    }
    pub fn dictionary_get_words_by_token(&self) {
        //exact-match
        //creates if not defined, including upsert
    }
    pub fn dictionary_get_words_by_id(&self) {
        //if not defined, raise an error
    }
    pub fn dictionary_set_words(&self) {
        //the level consuming this should delete case-sensitive variations of banned words
        //to save a bit of space
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
    pub fn load(&self, id:&str) -> Result<Database, Error> {
        let path = self.resolve_path(id);
        info!("loading database {}...", id);
        let connection = Connection::open(path)?;
        debug!("preparing database structures...");
        
        connection.execute("CREATE TABLE IF NOT EXISTS dictionary (
            caseInsensitiveRepresentation TEXT NOT NULL UNIQUE,
            id INTEGER NOT NULL PRIMARY KEY,
            caseInsensitiveOccurrences INTEGER NOT NULL,
            capitalisedFormsJSONZLIB BLOB
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
