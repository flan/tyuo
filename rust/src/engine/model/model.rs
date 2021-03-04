use crate::engine::model::database;

use fnv::{FnvHashMap, FnvHashSet};
use std::collections::HashMap;

use std::time::{SystemTime, UNIX_EPOCH};

const MAX_TRANSITION_AGE:i64 = 3600 * 24 * 365; //approximately one year
const DECIMATION_THRESHOLD:u16 = 100;
const DECIMATION_FACTOR:u8 = 3;

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
    
    pub fn increment(&mut self) -> i32 {
        self.occurrences += 1;
        self.last_observed = get_current_time();
        return self.occurrences.clone();
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
    transitions: HashMap<i32, Transition>,
}
impl Node {
    pub fn new(transitions:HashMap<i32, Transition>) -> Node {
        return Node{
            transitions: transitions,
        };
    }
    
    pub fn increment_transition(&mut self, dictionary_id:&i32) -> i32 {
        let transition = self.transitions.get_mut(dictionary_id);
        return match transition {
            Some(t) => t.increment(),
            None => {
                self.transitions.insert(
                    dictionary_id.to_owned(),
                    Transition::new(1, get_current_time()),
                );
                return 1;
            },
        };
    }
    
    //function to choose a transition... maybe
    //perhaps that should go in logic, unless it's sensible to defer that here
    //for encapsulation
}


pub struct Statistics {
    //database reference
    //direction identifier
    //banned words
}
impl Statistics {
    //function to get Nodes by dictionary IDs
    //when loading transitions, scrub banned targets and anything that has expired
    
    //function to save Nodes
}
