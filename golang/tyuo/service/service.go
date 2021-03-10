package service

import (
    ctx "context"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "time"
    
    "github.com/flan/tyuo/context"
)

var httpIp = flag.String("http-ip", "", "the IP on which to listen for HTTP requests")
var httpPort = flag.Uint("http-port", 8080, "the port on which to listen for HTTP requests")

var contextIdRe = regexp.MustCompile("^[_a-zA-Z0-9][-_a-zA-Z0-9]{0,220}$")


func doPreamble(w *http.ResponseWriter, r *http.Request) (bool) {
    if r.Method != http.MethodPost && r.Method != http.MethodOptions {
        http.Error(*w, fmt.Sprintf("unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
        return false
    }
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    if r.Method == http.MethodOptions {
        (*w).WriteHeader(http.StatusNoContent)
        return false
    }
    return true
}


type speakRequest struct {
    ContextId string
    Input string
}
func speakHandler(w http.ResponseWriter, r *http.Request, cm *context.ContextManager) {
    if !doPreamble(&w, r) {
        return
    }
    
    requestJson, err := ioutil.ReadAll(r.Body)
    if err != nil {
        logger.Errorf("unable to read HTTP body: %s", err)
        http.Error(w, "unable to read request", http.StatusInternalServerError)
        return
    }
    var request speakRequest
    if err = json.Unmarshal(requestJson, &request); err != nil {
        logger.Warningf("malformed JSON block: %s", requestJson)
        http.Error(w, "unable to parse JSON", http.StatusBadRequest)
        return
    }
    
    if !contextIdRe.MatchString(request.ContextId) {
        logger.Warningf("invalid context ID: %s", request.ContextId)
        http.Error(w, "invalid context ID", http.StatusBadRequest)
        return
    }
    /*
    context, err := cm.GetContext(request.ContextId)
    if err != nil {
        logger.Warningf("unable to access context %s: %s", request.ContextId, err)
        http.Error(w, "unable to access context", http.StatusBadRequest)
    }
    */
    
    var startTime time.Time = time.Now()
    
    /*
    logger.Infof("%d", context)
    //ask logic to generate a response, which internally calls Parse
    //and everything else in the chain
    */
    
    logger.Infof("prepared response in %s", time.Now().Sub(startTime))
}

func RunForever(shutdown chan<- string, contextManager *context.ContextManager) (chan<- bool) {
    var kill = make(chan bool, 1)
    var addr = fmt.Sprintf("%s:%d", *httpIp, *httpPort)
    
    srv := &http.Server{Addr: addr}
    
    go func() {
        http.HandleFunc("/speak", func(w http.ResponseWriter, r *http.Request) {
            speakHandler(w, r, contextManager)
        })
        /*
        http.HandleFunc("/learn", func(w http.ResponseWriter, r *http.Request) {
            learnHandler(w, r, contextManager)
        })
        
        http.HandleFunc("/ban", func(w http.ResponseWriter, r *http.Request) {
            banHandler(w, r, contextManager)
        })
        http.HandleFunc("/unban", func(w http.ResponseWriter, r *http.Request) {
            unbanHandler(w, r, contextManager)
        })
        */

        logger.Infof("starting HTTP service on %s...", addr)
        if err := srv.ListenAndServe(); err != nil {
            shutdown<- fmt.Sprintf("unable to serve HTTP: %s", err)
        }
    }()
    
    go func() {
        <-kill
        logger.Infof("shutting down HTTP service...")
        if err := srv.Shutdown(ctx.Background()); err != nil {
            logger.Errorf("unable to shut down HTTP service: %s", err)
        }
    }()
    
    return kill
}



//context-ID needs to be sanitised to make sure it isn't a path-spec.
//just make it match a-zA-Z0-9{1,220}

//use JSON to handle interactions
/*

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
