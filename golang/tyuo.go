/* tyuo is a Markov-chain-based chatbot, loosely based on MegaHAL by Jason
 * Hutchens.
 * 
 * More specifically, tyuo is a rewrite of yuo, which was cobbled together by
 * Neil Tallim in 2002, based on a limited understanding of how MegaHAL worked
 * and very rudimentary knowledge of C.
 * While things have probably diverged a lot, MegaHAL was undeniably the initial
 * inspiration for whatever this is now.
 */
package main

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

var dataDir = flag.String("data-dir", "", "the path to tyuo's data (default ~/.tyuo/)")

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
    
    var dataPath string
    if *dataDir == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            logger.Criticalf("unable to determine user details: %s", err)
            os.Exit(1)
        }
        dataPath = filepath.Join(homeDir, ".tyuo")
    } else {
        dataPath = *dataDir
    }
    
    
    fmt.Println(dataPath)
    language.Test(dataPath)
    
    
    contextManager, err := context.PrepareContextManager(dataPath)
    if err != nil {
        panic(err)
    }
    
    shutdownChannel := make(chan string, 1)
    
    setupSignals(shutdownChannel)
    
    httpShutdownChannel := service.RunForever(shutdownChannel, contextManager)
    
    logger.Infof("beginning normal operation...")
    logger.Warningf("system shutting down: %s...", <-shutdownChannel)
    
    httpShutdownChannel<- true
    
    contextManager.Close()
}
