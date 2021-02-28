use crate::engine::model::banned_dictionary;

use std::io::{Read, Write};

use rusqlite::{Connection, Error, params};
use serde_json::json;

fn to_json_zlib(serialisable_data:serde_json::Value) -> Vec<u8> {
    let mut encoder = flate2::write::ZlibEncoder::new(Vec::new(), flate2::Compression::default());
    encoder.write_all(&serde_json::to_vec(&serialisable_data).unwrap());
    return encoder.finish().unwrap();
}
fn from_json_zlib(blob:Vec<u8>) -> Option<serde_json::Value> {
    let mut decoder = flate2::read::ZlibDecoder::new(&blob[..]);
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
    
    
    pub fn banned_load_banned_tokens(&self) -> Result<Vec<banned_dictionary::BannedWord>, Error> {
        let mut stmt = self.connection.prepare("SELECT caseInsensitiveRepresentation FROM dictionary_banned")?;
        let iter = stmt.query_map(params![], |row| {
            let case_insensitive_representation:String = row.get(0)?;
            return Ok(banned_dictionary::BannedWord {
                case_insensitive_representation: case_insensitive_representation.to_lowercase(),
            });
        })?;
        
        let mut output = Vec::new();
        for banned_word in iter {
            output.push(banned_word.unwrap());
        }
        return Ok(output);
    }
    pub fn banned_ban_tokens(&self) {
        
    }
    pub fn banned_unban_tokens(&self) {
        
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
            caseInsensitiveRepresentation TEXT NOT NULL PRIMARY KEY,
        )", params![])?;
        connection.execute("CREATE TABLE IF NOT EXISTS statistics_forward (
            dictionaryId INTEGER NOT NULL PRIMARY KEY,
            childrenJSONZLIB BLOB,
            
            FOREIGN KEY(dictionaryId) REFERENCES dictionary(id) ON DELETE CASCADE,
        )", params![])?;
        connection.execute("CREATE TABLE IF NOT EXISTS statistics_reverse (
            dictionaryId INTEGER NOT NULL PRIMARY KEY,
            childrenJSONZLIB BLOB,
            
            FOREIGN KEY(dictionaryId) REFERENCES dictionary(id) ON DELETE CASCADE,
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
