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
    "github.com/flan/tyuo/logic"
)

var httpIp = flag.String("http-ip", "", "the IP on which to listen for HTTP requests (default all)")
var httpPort = flag.Uint("http-port", 48100, "the port on which to listen for HTTP requests")

var contextIdRe = regexp.MustCompile("^[_a-zA-Z0-9][-_a-zA-Z0-9]{0,220}$")


func doPreamble(w *http.ResponseWriter, r *http.Request) (*[]byte) {
    if r.Method != http.MethodPost && r.Method != http.MethodOptions {
        http.Error(*w, fmt.Sprintf("unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
        return nil
    }
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    if r.Method == http.MethodOptions {
        (*w).WriteHeader(http.StatusNoContent)
        return nil
    }
    
    requestJson, err := ioutil.ReadAll(r.Body)
    if err != nil {
        logger.Errorf("unable to read HTTP body: %s", err)
        http.Error(*w, "unable to read request", http.StatusInternalServerError)
        return nil
    }
    
    return &requestJson
}

func unmarshalRequest(w *http.ResponseWriter, r *http.Request, requestJson []byte, request interface{}) (error) {
    if err := json.Unmarshal(requestJson, request); err != nil {
        logger.Warningf("malformed JSON block: %s", requestJson)
        http.Error(*w, "unable to parse JSON", http.StatusBadRequest)
        return err
    }
    return nil
}

func getContext(w *http.ResponseWriter, r *http.Request, contextId string, cm *context.ContextManager) (*context.Context) {
    if !contextIdRe.MatchString(contextId) {
        logger.Warningf("invalid context ID: %s", contextId)
        http.Error(*w, "invalid context ID", http.StatusBadRequest)
        return nil
    }
    
    ctx, err := cm.GetContext(contextId)
    if err != nil {
        logger.Warningf("unable to access context %s: %s", contextId, err)
        http.Error(*w, "unable to access context", http.StatusBadRequest)
        return nil
    }
    return ctx
}



type speakRequest struct {
    ContextId string
    Input string
}
func speakHandler(w http.ResponseWriter, r *http.Request, cm *context.ContextManager) {
    requestJson := doPreamble(&w, r)
    if requestJson == nil {return}
    
    var request speakRequest
    if err := unmarshalRequest(&w, r, *requestJson, &request); err != nil {return}
    ctx := getContext(&w, r, request.ContextId, cm)
    if ctx == nil {return}
    
    
    var startTime time.Time = time.Now()
    
    logic.Speak(ctx, request.Input)
    
    logger.Infof("prepared response in %s", time.Now().Sub(startTime))
}

type learnRequest struct {
    ContextId string
    Input []string
}
func learnHandler(w http.ResponseWriter, r *http.Request, cm *context.ContextManager) {
    requestJson := doPreamble(&w, r)
    if requestJson == nil {return}
    
    var request learnRequest
    if err := unmarshalRequest(&w, r, *requestJson, &request); err != nil {return}
    ctx := getContext(&w, r, request.ContextId, cm)
    if ctx == nil {return}
    
    
    var startTime time.Time = time.Now()
    
    linesLearned := logic.Learn(ctx, request.Input)
    
    logger.Infof("learned %d lines of input in %s", linesLearned, time.Now().Sub(startTime))
}

type banRequest struct {
    ContextId string
    Substrings []string
}
func banSubstringsHandler(w http.ResponseWriter, r *http.Request, cm *context.ContextManager) {
    requestJson := doPreamble(&w, r)
    if requestJson == nil {return}
    
    var request banRequest
    if err := unmarshalRequest(&w, r, *requestJson, &request); err != nil {return}
    ctx := getContext(&w, r, request.ContextId, cm)
    if ctx == nil {return}
    
    
    var startTime time.Time = time.Now()
    
    logic.BanSubstrings(ctx, request.Substrings)
    
    logger.Infof("banned in %s", time.Now().Sub(startTime))
}
func unbanSubstringsHandler(w http.ResponseWriter, r *http.Request, cm *context.ContextManager) {
    requestJson := doPreamble(&w, r)
    if requestJson == nil {return}
    
    var request banRequest
    if err := unmarshalRequest(&w, r, *requestJson, &request); err != nil {return}
    ctx := getContext(&w, r, request.ContextId, cm)
    if ctx == nil {return}
    
    
    var startTime time.Time = time.Now()
    
    logic.UnbanSubstrings(ctx, request.Substrings)
    
    logger.Infof("unbanned in %s", time.Now().Sub(startTime))
}


func RunForever(shutdown chan<- string, contextManager *context.ContextManager) (chan<- bool) {
    var kill = make(chan bool, 1)
    var addr = fmt.Sprintf("%s:%d", *httpIp, *httpPort)
    
    srv := &http.Server{Addr: addr}
    
    go func() {
        http.HandleFunc("/speak", func(w http.ResponseWriter, r *http.Request) {
            speakHandler(w, r, contextManager)
        })
        http.HandleFunc("/learn", func(w http.ResponseWriter, r *http.Request) {
            learnHandler(w, r, contextManager)
        })
        
        http.HandleFunc("/banSubstrings", func(w http.ResponseWriter, r *http.Request) {
            banSubstringsHandler(w, r, contextManager)
        })
        http.HandleFunc("/unbanSubstrings", func(w http.ResponseWriter, r *http.Request) {
            unbanSubstringsHandler(w, r, contextManager)
        })

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
