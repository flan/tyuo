mod logic;
mod model;

pub fn goodbye() {
    logic::goodbye();
}

pub struct Context {
    //dictionary
    //banned dictionary (pointer to engine's)
    
    //model (includes backwards and forwards paths)
}
impl Context {
    pub fn test(&self) -> &str {
        return "hello";
    }
}

pub struct Engine {
    //banned dictionary
    
    //map of established contexts
}
impl Engine {
    pub fn get_context(&self, id:&str) -> Result<Context, String> {
        return Ok(Context{
        });
    }
}

pub fn prepare(memory_dir:&std::path::Path, banned_tokens_list:&std::path::Path) -> Result<Engine, String> {
    return Ok(Engine{
    });
}
