#[macro_use]
extern crate log;

mod engine;
mod service;

fn main() {
    env_logger::builder()
        .filter(None, log::LevelFilter::Info)
        .init();
        
    let matches = clap::App::new("tyuo")
        .version("0.0.1")
        .author("Neil Tallim <flan@uguu.ca>")
        .about("Markov-chain-based chatter action")
        .arg(clap::Arg::new("memory-dir")
            .long("memory-dir")
            .about("the path in which tyuo's memories are stored")
            .default_value(dirs::home_dir().unwrap().join(".tyuo/memories").to_str().unwrap())
            .takes_value(true))
        .arg(clap::Arg::new("banned-tokens-list")
            .long("banned-tokens-list")
            .about("the path to a file containing banned tokens")
            .default_value(dirs::home_dir().unwrap().join(".tyuo/banned-tokens").to_str().unwrap())
            .takes_value(true))
        .get_matches();
        
    service::hello();
    engine::goodbye();
    
    info!("Preparing engine...");
    let engine = engine::prepare(
        std::path::Path::new(matches.value_of("memory-dir").unwrap()),
        std::path::Path::new(matches.value_of("banned-tokens-list").unwrap()),
    ).unwrap();
    
    info!("Preparing service...");
    //let service = service::prepare(engine);
    
    info!("Running service...");
    //service.serve_forever();
    
    error!("{}", engine.get_context("whee").unwrap().test());
    
    //engine.shutdown();
}
