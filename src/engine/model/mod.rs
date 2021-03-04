mod banned_dictionary;
mod database;
mod dictionary;
mod model;

use std::error::Error;

pub fn goodbye() {
    println!("Goodbye, world!");
}

pub struct Model<'md> {
    database_manager: database::DatabaseManager<'md>,
    
    banned_tokens_generic: Vec<String>,
    //non-keyword list (origin)
    
    //contexts {id: {model(database, dictionary_banned), dictionary(database, non-keyword-tokens list, dictionary_banned), dictionary_banned(database, banned list)}}
}
impl<'md> Model<'md> {
    pub fn new(
        db_dir:&std::path::Path,
        non_keyword_tokens:&'md std::path::Path,
        banned_tokens_list:&'md std::path::Path,
        parsing_language:&str,
    ) -> Result<Model<'md>, Box<Error>> {
        
        
        //TODO: dev test
        /*let dbr = dbm.load("hi");
        if dbr.is_err(){
            eprintln!("{:?}", dbr.err());
        } else {
            let mut db = dbr.unwrap();
            //these use HashSets now
            println!("{:?}", db.banned_ban_tokens(vec!["hi", "bye", "test"]).unwrap());
            //after calling this, iterate over anything that comes back
            //and, if it has a dictionary ID, delete related transitions
            //and scrub capitalised forms of dictionary entries, setting
            //the insensitive count to 0
            
            println!("{:?}", db.banned_load_banned_tokens(Some(vec!["hi", "bye", "test"])).unwrap());
            println!("{:?}", db.banned_unban_tokens(vec!["bye"]));
            println!("{:?}", db.banned_load_banned_tokens(None).unwrap());
        }
        */
        
        let database_manager = database::DatabaseManager::new(db_dir);
        
        let banned_tokens_generic = banned_dictionary::process_banned_tokens_list(
            banned_tokens_list,
        )?;
        
        println!("{:?}", banned_tokens_generic);
        
        let bd = banned_dictionary::BannedDictionary::new(
            &mut Box::new(database_manager.load("hi").unwrap()),
            &banned_tokens_generic,
        );
        bd.ban(&vec!["hello", "desu", "bye"])?;
        println!("{:?}", db.is_banned_by_token(&vec!["goodbye", "oh"]));
        bd.unban(&vec!["bye"])?;
        println!("{:?}", db.is_banned_by_token(&vec!["goodbye", "desu"]));
        
        println!("hi2");
        
        return Ok(Model{
            database_manager: database_manager,
            
            banned_tokens_generic: banned_tokens_generic,
        });
    }

    pub fn get_context(&mut self, id:&str) {

    }
    pub fn create_context(&self, id:&str) {

    }
    pub fn drop_context(&mut self, id:&str) {

    }
}
