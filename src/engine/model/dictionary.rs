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
            case_insensitive_representation: case_insensitive_representation.to_lowercase(),
            capitalised_forms: capitalised_forms,
        };
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
}
impl Dictionary {
    //function to get a collection of words as a DictionarySlice, holding them
    //in memory, keyed by ID and case-insensitive repr
    //this is used to efficiently perform updates while learning and format output
    //two forms to get this, actually: one by IDs and one by tokens
}
