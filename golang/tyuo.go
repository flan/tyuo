package main

//input is just "dataDir"; other paths are assumed relative to that, based on language
//when creating a context, the language will need to be specified
//it is then loaded afterwards, alongside the database file

import (
    "fmt"
    
    "os"
    "path/filepath"
    
    "github.com/juju/loggo"
    
    "github.com/flan/tyuo/context"
)

var logger = loggo.GetLogger("main")

func main() {
    writer, _ := loggo.RemoveWriter("default")
    loggo.RegisterWriter("default", writer)
    
    homeDir, err := os.UserHomeDir()
    if err != nil {
        logger.Criticalf("unable to determine user details: %s", err)
        os.Exit(1)
    }
    
    contextDir := filepath.Join(homeDir, ".tyuo/contexts")
    fmt.Println(contextDir)
    context.Test(contextDir)
}
