use crate::engine::model::database;

use std::error::Error;
use std::fs::File;
use std::io::{BufRead, BufReader};

use std::collections::HashSet;
use fnv::FnvHashMap;

#[derive(Clone)]
pub struct BannedToken {
    case_insensitive_representation: String,
    dictionary_id: Option<i32>,
}
impl BannedToken {
    pub fn new(
        case_insensitive_representation:String,
        dictionary_id: Option<i32>,
    ) -> BannedToken {
        return BannedToken{
            case_insensitive_representation: case_insensitive_representation,
            dictionary_id: dictionary_id,
        };
    }
    
    pub fn token_equal(&self, token:&str) -> bool {
        return self.case_insensitive_representation == token;
    }
    pub fn token_contains(&self, token:&str) -> bool {
        return token.contains(self.case_insensitive_representation.as_str());
    }
    
    pub fn get_dictionary_id(&self) -> Option<i32> {
        return match self.dictionary_id {
            Some(v) => Some(v.clone()),
            None => None,
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
) -> Result<Vec<String>, Box<Error>> {
    let file = File::open(banned_tokens_list)?;
    let reader = BufReader::new(file);
    
    let mut tokens = Vec::new();
    for line in reader.lines() {
        let line = (line?).to_lowercase();
        let trimmed_line = line.trim();
        if !trimmed_line.is_empty() {
            tokens.push(trimmed_line.to_owned());
        }
    }
    return Ok(tokens);
}

pub struct BannedDictionary<'bd> {
    database: &'bd mut Box<database::Database>,
    
    //words from database
    banned_tokens: Vec<BannedToken>,
    banned_tokens_by_id: FnvHashMap<i32, BannedToken>,
    
    //tokens from the list
    banned_tokens_generic: &'bd Vec<String>,
}
impl<'bd> BannedDictionary<'bd> {
    pub fn new(
        database:&'bd mut Box<database::Database>,
        banned_tokens_generic:&'bd Vec<String>,
    ) -> Result<BannedDictionary<'bd>, Box<Error>> {
        let mut banned_tokens = Vec::new();
        let mut banned_tokens_by_id = FnvHashMap::default();
        for banned_token in database.banned_load_banned_tokens(None)? {
            banned_tokens.push(banned_token.clone());
            
            let dictionary_id = banned_token.get_dictionary_id();
            if dictionary_id.is_some() {
                banned_tokens_by_id.insert(
                    banned_token.get_dictionary_id().unwrap(),
                    banned_token,
                );
            }
        }
        
        return Ok(BannedDictionary{
            database: database,
            
            banned_tokens: banned_tokens,
            banned_tokens_by_id: banned_tokens_by_id,
            
            banned_tokens_generic: banned_tokens_generic,
        });
    }
    
    pub fn ban(&'bd mut self,
        tokens:&Vec<&str>,
    ) -> Result<(), Box<Error>> {
        let mut banned_tokens = HashSet::new();
        for token in tokens {
            let mut already_banned = false;
            for banned_token in &self.banned_tokens {
                if banned_token.token_equal(token) {
                    already_banned = true;
                    break;
                }
            }
            if !already_banned {
                banned_tokens.insert(*token);
            }
        }
        
        let newly_banned_tokens = self.database.banned_ban_tokens(&banned_tokens)?;
        for banned_token in newly_banned_tokens {
            let bt = banned_token.clone();
            self.banned_tokens.push(bt.clone());
            let dictionary_id = banned_token.get_dictionary_id();
            if dictionary_id.is_some() {
                self.banned_tokens_by_id.insert(
                    dictionary_id.unwrap(),
                    banned_token,
                );
            }
        }
        return Ok(());
    }
    
    pub fn unban(&'bd mut self,
        tokens:&Vec<&str>,
    ) -> Result<(), Box<Error>> {
        let mut unbanned_tokens = HashSet::new();
        for token in tokens {
            for banned_token in &self.banned_tokens {
                if banned_token.token_equal(token) {
                    unbanned_tokens.insert(*token);
                    break;
                }
            }
        }
        
        self.database.banned_unban_tokens(&unbanned_tokens)?;
        
        let mut banned_tokens_by_id = &mut self.banned_tokens_by_id;
        self.banned_tokens.retain(|bt| {
            for unbanned_token in &unbanned_tokens {
                if bt.token_equal(unbanned_token) {
                    match bt.get_dictionary_id() {
                        Some(v) => {
                            banned_tokens_by_id.remove(&v);
                            ()
                        },
                        None => (),
                    };
                    return false;
                }
            }
            return true;
        });
        
        return Ok(());
    }
    
    pub fn is_banned_by_token(&self, tokens:&Vec<&str>) -> bool {
        for banned_token in &self.banned_tokens {
            for token in tokens {
                if banned_token.token_contains(token) {
                    return true;
                }
            }
        }
        for banned_token in self.banned_tokens_generic {
            for token in tokens {
                if token.contains(banned_token) {
                    return true;
                }
            }
        }
        return false;
    }
    
    pub fn is_banned_by_id(&self, ids:&Vec<i32>) -> bool {
        for id in ids {
            if self.banned_tokens_by_id.contains_key(&id) {
                return true;
            }
        }
        return false;
    }
}
