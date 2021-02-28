use crate::engine::model::database;

pub struct BannedWord {
    case_insensitive_representation: String,
    dictionary_id: Option<i32>,
}
impl BannedWord {
    pub fn prepare(case_insensitive_representation:String, dictionary_id: Option<i32>) -> BannedWord {
        return BannedWord{
            case_insensitive_representation: case_insensitive_representation.to_lowercase(),
            dictionary_id: dictionary_id,
        };
    }
}
