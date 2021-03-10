//context-ID needs to be sanitised to make sure it isn't a path-spec.
//just make it match a-zA-Z0-9{1,220}


//run a TCP server to handle interactions

//use JSON to handle interactions
/*
 {
    "action": "prompt",
    "context": <ID as string>,
    "input": [<string>],
    "learn": <bool>, //should this be an option? Maybe just explicitly requiring a second call will be clearer
 }
 {
    "action": "learn",
    "context": <ID as string>,
    "input": [<string>],
 }
 {
    "action": "banTokens",
    "context": <ID as string>,
    "tokens": [<token>],
 }
 {
    "action": "unbanTokens",
    "context": <ID as string>,
    "tokens": [<token>],
 }
*/
