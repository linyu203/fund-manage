
package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    //"github.com/gorilla/mux"
)

var (
    // See template.go.
    indexTmpl    *template.Template//= "index.html"
    bondListTmpl *template.Template//= "BondList.html"
    newFundTmpl  *template.Template//= "NewFund.html"
)

func initTemplates() {
    var err error = nil
    indexTmpl, err = template.ParseFiles("templates/index.html")
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
    http.HandleFunc("/", getAllFundHandler)
    http.HandleFunc("/bonds/{fundName:[0-9a-zA-Z_\\-]+}", bondsHandler)



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
    fmt.Printf("r: %#v", &r)
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
    switch r.Method {
    case "GET":
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
    default:
        http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
    }
}



func main(){
    initTemplates()
    InitTable()
    regestHandler()
}
