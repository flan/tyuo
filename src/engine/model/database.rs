pub fn oh_no() {
    info!("oh no");
}

pub struct Database {
    
}
impl Database {
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
    pub fn create(&self, id:&str) { //TODO: may fail
        
    }
    pub fn drop(&self, id:&str) { //TODO: may fail
        
    }
}

