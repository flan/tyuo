use crate::engine::model::database;

pub struct BannedWord {
    pub case_insensitive_representation: String,
    pub dictionary_id: Option<u32>,
}
