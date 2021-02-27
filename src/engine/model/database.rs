pub fn oh_no() {
    info!("oh no");
}

use rusqlite::{Connection, Error, params};

pub struct Database {
    connection: Connection,
}
impl Database {
    pub fn prepare(connection:Connection) -> Database {
        return Database{
            connection: connection,
        };
    }
    
    //functions to read and write the dictionary
    
    //functions to read and write statistics
    
    
    pub fn test(&self) -> &str {
        return "hello";
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
