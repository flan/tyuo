mod banned_dictionary;
mod database;
mod dictionary;

pub fn goodbye() {
    println!("Goodbye, world!");
}

pub struct Model {
    database_manager: Box<database::DatabaseManager>,
    //banned list (origin)
    //non-keyword list (origin)
    
    //contexts {id: {model(database, dictionary_banned), dictionary(database, non-keyword-tokens list, dictionary_banned), dictionary_banned(database, banned list)}}
}
impl Model {
    pub fn prepare(
        db_dir:&std::path::Path,
        non_keyword_tokens:&std::path::Path,
        banned_tokens_list:&std::path::Path,
    ) -> Model {
        //TODO: dev test
        let dbm = database::DatabaseManager::prepare(db_dir);
        let dbr = dbm.load("hi");
        if dbr.is_err(){
            eprintln!("{:?}", dbr.err());
        } else {
            let mut db = dbr.unwrap();
            /* //these use HashSets now
            println!("{:?}", db.banned_ban_tokens(vec!["hi", "bye", "test"]).unwrap());
            //after calling this, iterate over anything that comes back
            //and, if it has a dictionary ID, delete related transitions
            //and scrub capitalised forms of dictionary entries, setting
            //the insensitive count to 0
            
            println!("{:?}", db.banned_load_banned_tokens(Some(vec!["hi", "bye", "test"])).unwrap());
            println!("{:?}", db.banned_unban_tokens(vec!["bye"]));
            println!("{:?}", db.banned_load_banned_tokens(None).unwrap());
            */
        }

        return Model{
            database_manager: database::DatabaseManager::prepare(db_dir),
        };
    }

    pub fn get_context(&mut self, id:&str) {

    }
    pub fn create_context(&self, id:&str) {

    }
    pub fn drop_context(&mut self, id:&str) {

    }
}
