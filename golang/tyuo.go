/* tyuo is a Markov-chain-based chatbot, loosely based on MegaHAL by Jason
 * Hutchens.
 * 
 * More specifically, tyuo is a rewrite of yuo, written by Neil Tallim in 2002,
 * based on a limited understanding of how MegaHAL worked, significantly
 * butchered, but that was undeniably the initial inspiration for whatever this
 * is now.
 */
package main

//input is just "dataDir"; other paths are assumed relative to that, based on language
//when creating a context, the language will need to be specified
//it is then loaded afterwards, alongside the database file

import (
    "fmt"
    "flag"
    
    "os"
    "os/signal"
    "path/filepath"
    
    "syscall"
    
    "github.com/juju/loggo"
    
    "github.com/flan/tyuo/context"
    "github.com/flan/tyuo/logic/language"
    "github.com/flan/tyuo/service"
)

var logger = loggo.GetLogger("main")

func setupSignals(shutdown chan<- string) {
    var sigs = make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        var sig = <-sigs
        shutdown<- fmt.Sprintf("requested by operator, signal=%s", sig.String())
    }()
}

func main() {
    flag.Parse()
    //use a flag to set logging level
    
    loggo.GetLogger("").SetLogLevel(loggo.DEBUG)
    loggo.RemoveWriter("default")
    loggo.RegisterWriter("console", loggo.NewSimpleWriter(os.Stderr, loggo.DefaultFormatter))
    
    
    homeDir, err := os.UserHomeDir()
    if err != nil {
        logger.Criticalf("unable to determine user details: %s", err)
        os.Exit(1)
    }
    
    contextDir := filepath.Join(homeDir, ".tyuo/contexts")
    fmt.Println(contextDir)
    context.Test(contextDir)
    language.Test(contextDir)
    
    shutdownChannel := make(chan string, 1)
    
    setupSignals(shutdownChannel)
    
    httpShutdownChannel := service.RunForever(shutdownChannel, nil)
    
    logger.Infof("beginning normal operation...")
    logger.Warningf("system shutting down: %s...", <-shutdownChannel)
    
    httpShutdownChannel<- true
}
