
package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "io"
	"os"
    "github.com/gorilla/mux"
)

var (
    // See template.go.
    indexTmpl    *template.Template//= "index.html"
    bondListTmpl *template.Template//= "BondList.html"
    newFundTmpl  *template.Template//= "NewFund.html"
    logWriter    io.Writer
)

func initTemplates() {
    var err error = nil
    indexTmpl, err = template.ParseFiles("templates/index.html")
    logWriter = os.Stderr
    if err == nil {
        bondListTmpl, err = template.ParseFiles("templates/BondList.html")
    }
    if err == nil {
        //newFundTmpl, err = template.ParseFiles("templates/NewFund.html")
    }
    if err != nil {
        log.Fatalf("unable to parse template file: %s", err)
    }
}

func regestHandler(){
    r := mux.NewRouter()
    r.Methods("GET").Path("/").Handler(appHandler(getAllFundHandler))
    r.Methods("GET").Path("/books/{id:[0-9a-zA-Z_\\-]+}").Handler(appHandler(bondsHandler))
    http.Handle("/", handlers.CombinedLoggingHandler(logWriter, r))

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Listening on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}

func bondsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("bondsHandler called %v", r.Method)
    /*switch r.Method {
    case "GET":
        bonds, err := GetAllfunds()
        if err != nil {
            log.Printf("Get all funds: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        err = indexTmpl.Execute(w, funds)
        fmt.Printf("funds: %#v\n", funds)
        if err != nil {
            log.Printf("Transefor to html error: %v", err)
        }
    default:
        http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
    }*/
}
func getAllFundHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("getAllFundHandler called %v", r.Method)
    funds, err := GetAllfunds()
    if err != nil {
        log.Printf("Get all funds: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    err = indexTmpl.Execute(w, funds)
    fmt.Printf("funds: %#v\n", funds)
    if err != nil {
        log.Printf("Transefor to html error: %v", err)
    }
}



func main(){
    initTemplates()
    InitTable()
    regestHandler()
}
