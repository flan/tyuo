use crate::engine::model::database;

use std::collections::HashMap;

pub struct DictionaryWord {
    id: i32,
    case_insensitive_occurrences: i32,
    case_insensitive_representation: String,
    capitalised_forms: HashMap<String, i32>,
}
impl DictionaryWord {
    pub fn prepare(
        id:i32,
        case_insensitive_occurrences:i32,
        case_insensitive_representation:String,
        capitalised_forms: HashMap<String, i32>,
    ) -> DictionaryWord {
        return DictionaryWord{
            id: id,
            case_insensitive_occurrences: case_insensitive_occurrences,
            case_insensitive_representation: case_insensitive_representation,
            capitalised_forms: capitalised_forms,
        };
    }
    
    pub fn get_id(&self) -> i32 {
        return self.id.clone();
    }
    pub fn get_case_insensitive_occurrences(&self) -> i32 {
        return self.case_insensitive_occurrences.clone();
    }
    pub fn get_case_insensitive_representation(&self) -> String {
        return self.case_insensitive_representation.clone();
    }
    pub fn get_capitalised_forms(&self) -> HashMap<String, i32> {
        return self.capitalised_forms.clone();
    }
    
    pub fn represent(&self, firstToken:bool) -> String {
        //see Go for logic
        return "".to_string();
    }
}
impl std::fmt::Debug for DictionaryWord {
    fn fmt(&self, fmt: &mut std::fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(fmt, "DictionaryWord {{ repr: {:?}, id: {:?}, occ: {:?}, forms: {:?} }}",
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
    //a list of all words that are never keyword candidates
}
impl Dictionary {
    //functions to get a collection of words by token or by ID
    
    //function to derive a HashSet of IDs from a given input string
    //if too few words qualify, a random sample of keyword candidates
    //is added in a loop, until the threshold is satisfied.
    //there should probably be a priority ordering of primary candidates (from input)
    //and then secondaries chosen at random
    
    //function to learn from a given input-string, updating the dictionary and
    //returning the identifiers of all tokens received, in sequential order
}
