use crate::engine::model::database;

use std::collections::HashMap;

pub struct DictionaryToken {
    id: i32,
    case_insensitive_occurrences: i32,
    case_insensitive_representation: String,
    capitalised_forms: HashMap<String, i32>,
}
impl DictionaryToken {
    pub fn new(
        id:i32,
        case_insensitive_occurrences:i32,
        case_insensitive_representation:String,
        capitalised_forms: HashMap<String, i32>,
    ) -> DictionaryToken {
        return DictionaryToken{
            id: id,
            case_insensitive_occurrences: case_insensitive_occurrences,
            case_insensitive_representation: case_insensitive_representation,
            capitalised_forms: capitalised_forms,
        };
    }
    
    pub fn get_id(&self) -> &i32 {
        return &self.id;
    }
    pub fn get_case_insensitive_occurrences(&self) -> &i32 {
        return &self.case_insensitive_occurrences;
    }
    pub fn get_case_insensitive_representation(&self) -> &String {
        return &self.case_insensitive_representation;
    }
    pub fn get_capitalised_forms(&self) -> &HashMap<String, i32> {
        return &self.capitalised_forms;
    }
    
    pub fn represent(&self, first_token:bool) -> String {
        //see Go for logic
        return "".to_string();
    }
}
impl std::fmt::Debug for DictionaryToken {
    fn fmt(&self, fmt: &mut std::fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(fmt, "DictionaryToken {{ repr: {:?}, id: {:?}, occ: {:?}, forms: {:?} }}",
            self.case_insensitive_representation,
            self.id,
            self.case_insensitive_occurrences,
            self.capitalised_forms,
        )?;
        Ok(())
    }
}

pub struct Dictionary {
    //database reference
    //latest ID
    //a list of all tokens that are never keytoken candidates
    //banned tokens
}
impl Dictionary {
    //functions to get a collection of tokens by token or by ID
    
    //function to derive a HashSet of IDs from a given input string as keytokens
    //if too few tokens qualify, a random sample of keytoken candidates
    //is added in a loop, until the threshold is satisfied.
    //there should probably be a priority ordering of primary candidates (from input)
    //and then secondaries chosen at random
    //the random set may contain banned tokens, so do an ID check before allowing
    //them to become candidates
    
    //function to learn from a given input-string, updating the dictionary and
    //returning the identifiers of all tokens received, in sequential order
}

//TODO: There should be another module that's entirely focused on tokenising
//input before it reaches the dictionary, to ensure punctuation consistency
//(and also to filter out banned content before learning can occur)
//this thing will return a vector of allowed tokens, both for learning and
//responding, and a boolean value indicating whether learning is even possible.
//It should be the same input-function for all cases, though.
