use crate::engine::model::database;

use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};

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
) -> Result<Vec<Box<str>>, Box<Error>> {
    let file = File::open(banned_tokens_list)?;
    let reader = BufReader::new(file);
    
    let mut tokens = Vec::new();
    for line in reader.lines() {
        let line = (line?).to_lowercase();
        let trimmed_line = line.trim();
        if !trimmed_line.is_empty() {
            tokens.push(trimmed_line.to_owned().into_boxed_str());
        }
    }
    return Ok(tokens);
}

pub struct BannedDictionary<'bd> {
    database: &'bd mut database::Database,
    
    banned_tokens: Vec<BannedToken>, //words from database
    banned_tokens_generic: &'bd Vec<Box<str>>, //tokens from the list
}
impl<'bd> BannedDictionary<'bd> {
    pub fn new(
        database:&'bd mut database::Database,
        banned_tokens_generic:&'bd Vec<Box<str>>,
    ) -> Result<BannedDictionary<'bd>, Box<Error>> {
        let banned_tokens = database.banned_load_banned_tokens(None)?;
        return Ok(BannedDictionary{
            database: database,
            
            banned_tokens: banned_tokens,
            banned_tokens_generic: banned_tokens_generic,
        });
    }
}
