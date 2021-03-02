use crate::engine::model::database;

pub struct BannedWord {
    case_insensitive_representation: String,
    dictionary_id: Option<i32>,
}
impl BannedWord {
    pub fn new(case_insensitive_representation:String, dictionary_id: Option<i32>) -> BannedWord {
        return BannedWord{
            case_insensitive_representation: case_insensitive_representation,
            dictionary_id: dictionary_id,
        };
    }
}
impl std::fmt::Debug for BannedWord {
    fn fmt(&self, fmt: &mut std::fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(fmt, "BannedWord {{ repr: {:?}, id: {:?} }}",
            self.case_insensitive_representation,
            self.dictionary_id,
        )?;
        Ok(())
    }
}
