use crate::engine::model::database;

use fnv::{FnvHashMap, FnvHashSet};

use std::time::{SystemTime, UNIX_EPOCH};

const MAX_TRANSITION_AGE:i64 = 3600 * 24 * 365; //approximately one year

fn get_current_time() -> i64 {
    let start = SystemTime::now();
    let since_the_epoch = start.duration_since(UNIX_EPOCH).unwrap();
    return since_the_epoch.as_secs() as i64;
}

pub struct Transition {
    occurrences: i32,
    last_observed: i64,
}
impl Transition {
    pub fn new(
        occurrences:i32,
        last_observed:i64,
    ) -> Transition {
        return Transition{
            occurrences: occurrences,
            last_observed: last_observed,
        };
    }
    
    pub fn get_occurrences(&self) -> &i32 {
        return &self.occurrences;
    }
    pub fn get_last_observed(&self) -> &i64 {
        return &self.last_observed;
    }
}
impl std::fmt::Debug for Transition {
    fn fmt(&self, fmt: &mut std::fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(fmt, "Transition {{ occ: {:?}, last: {:?} }}",
            self.occurrences,
            self.last_observed,
        )?;
        Ok(())
    }
}

pub struct Node {
    //FnvHashMap of transitions, keyed by dictionary ID
}
impl Node {
    //function to increment (or create at 1) transitions within
    
    //function to choose a transition
}


pub struct Model {
    //database reference
    //direction identifier
    //banned words
}
impl Model {
    //function to get Nodes by dictionary IDs
    //when loading transitions, scrub banned targets and anything that has expired
    
    //function to save Nodes
}
