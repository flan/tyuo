#[macro_use]
extern crate log;

mod logic;
mod model;
mod service;

fn main() {
    env_logger::builder()
        .filter(None, log::LevelFilter::Info)
        .init();
        
    service::hello();
    logic::goodbye();
    
    
    
    info!("starting up");
}
