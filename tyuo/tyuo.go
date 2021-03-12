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
    "gopkg.in/natefinch/lumberjack.v2"
    
    "github.com/flan/tyuo/context"
    "github.com/flan/tyuo/service"
)

var dataDir = flag.String("data-dir", "", "the path to tyuo's data (default ~/.tyuo/)")
var logLevel = flag.String("log-level", "warn", "the logging-level to be used")
var logConsole = flag.Bool("log-console", false, "whether to log to console")
var logFile = flag.String("log-file", "", "the path to which logs should be written")

var logger = loggo.GetLogger("main")

func setupSignals(shutdown chan<- string) {
    var sigs = make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        var sig = <-sigs
        shutdown<- fmt.Sprintf("requested by operator, signal=%s", sig.String())
    }()
}

func setupLogging() {
    var logLevelEnumerated, logLevelValid = loggo.ParseLevel(*logLevel)
    if !logLevelValid {
        logger.Errorf("Unsupported logging-level", *logLevel)
        logLevelEnumerated = loggo.WARNING
    }
    loggo.GetLogger("").SetLogLevel(logLevelEnumerated)
    
    loggo.RemoveWriter("default")
    
    if *logConsole {
        var writer = loggo.NewSimpleWriter(os.Stderr, loggo.DefaultFormatter)
        loggo.RegisterWriter("console", writer)
    }
    
    if *logFile != "" {
        var writer = loggo.NewSimpleWriter(&lumberjack.Logger{
            Filename: *logFile,
            MaxSize: 1, //megabytes
            MaxBackups: 3,
            MaxAge: 7, //days to hold backups
        }, nil)
        loggo.RegisterWriter("file", writer)
    }
}

func main() {
    flag.Parse()
    setupLogging()
    
    var dataPath string
    if *dataDir == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            panic(fmt.Sprintf("unable to determine user details: %s", err))
        }
        dataPath = filepath.Join(homeDir, ".tyuo")
    } else {
        dataPath = *dataDir
    }
    
    
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
