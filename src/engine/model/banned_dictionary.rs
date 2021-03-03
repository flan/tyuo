use crate::engine::model::database;

use std::error::Error;

pub struct BannedToken {
    case_insensitive_representation: String,
    dictionary_id: Option<i32>,
}
impl BannedToken {
    pub fn new(case_insensitive_representation:String, dictionary_id: Option<i32>) -> BannedToken {
        return BannedToken{
            case_insensitive_representation: case_insensitive_representation,
            dictionary_id: dictionary_id,
        };
    }
}
impl std::fmt::Debug for BannedToken {
    fn fmt(&self, fmt: &mut std::fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(fmt, "BannedToken {{ repr: {:?}, id: {:?} }}",
            self.case_insensitive_representation,
            self.dictionary_id,
        )?;
        Ok(())
    }
}

pub fn process_banned_tokens_list(
    banned_tokens_list:&std::path::Path,
) -> Vec<String> {
    return Vec::new();
}

pub struct BannedDictionary<'a> {
    database: &'a mut database::Database,
    
    banned_tokens: Vec<BannedToken>, //words from database
    banned_tokens_generic: &'a Vec<String>, //tokens from the list
}
impl<'a> BannedDictionary<'a> {
    pub fn new(
        database:&'a mut database::Database,
        banned_tokens_generic:&'a Vec<String>,
    ) -> Result<BannedDictionary<'a>, Box<Error>> {
        let banned_tokens = database.banned_load_banned_tokens(None)?;
        return Ok(BannedDictionary{
            database: database,
            
            banned_tokens: banned_tokens,
            banned_tokens_generic: banned_tokens_generic,
        });
    }
}
