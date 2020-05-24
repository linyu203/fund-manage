
package main

import (
    "fmt"
    "html/template"
    "bytes"
    "log"
    "io"
    "net/http"
    "os"
    "github.com/gorilla/mux"
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
        newFundTmpl, err = template.ParseFiles("templates/NewFund.html")
    }
    if err != nil {
        log.Fatalf("unable to parse template file: %s", err)
    }
}

func regestHandler(){
    r := mux.NewRouter()
    r.Path("/").Handler(appHandler(getAllFundHandler))
    r.Path("/bonds/{fund:[0-9a-zA-Z_\\- ]+}").Handler(appHandler(bondsHandler))
    r.HandleFunc("/fund", fundHandler)

    http.Handle("/",r)
    port := os.Getenv("PORT")
    if port == "" {
        port = ":8080"
    }
    log.Printf("Listening on port %s", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal(err)
    }
}
type FundDetail struct {
    FundName string
    Description string
    Bonds []Bond
}

type NewFund struct {
    Action string
    Name string
    Description string
}

type appHandler func(http.ResponseWriter, *http.Request)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fn(w,r)
}

func fundHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("fundHandler called method %v\n", r.Method)
    switch r.Method {
    case "GET":
	    err := newFundTmpl.Execute(w, nil)
        //fmt.Printf("funds: %#v\n", funds)
        if err != nil {
            log.Printf("Transefor to newFundTmpl html error: %v", err)
        }
    case "POST":
	defer func (){
            http.Redirect(w, r, "/", http.StatusFound)
	}()
	var name, desc string
	mr,err := r.MultipartReader()
	if err != nil {
	     fmt.Println("multi read faild", err)
	}else{
		for part,err := mr.NextPart(); err!=io.EOF; part,err = mr.NextPart() {
			//fmt.Println(part)
			if part.FormName() == "name" {
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				name = buf.String()
			}else if part.FormName() == "description"{
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)
				desc = buf.String()
			}
		}
	}


        //fmt.Printf("Create New Fund name %s des %s \n", name, desc)
	if name == "" || desc == "" {
		log.Printf("name of decription is not valid")
	}
        if err:= CreateFund(name, desc); err != nil {
            log.Printf("Create fund error: %v", err)
            //http.Error(w, "Create fund error", http.StatusInternalServerError)
        }else{
	    log.Printf("Fund %s created Successfully", name)
            //http.Error(w, "Fund created Successfully!", http.StatusOK)

        }

    default:
        http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
    }
}

func bondsHandler(w http.ResponseWriter, r *http.Request) {
    fundName := mux.Vars(r)["fund"]
    fmt.Printf("bondsHandler called method %v %s \n", r.Method, fundName)
    
    switch r.Method {
    case "GET":
        bonds, err := Getbonds(fundName)
        if err != nil {
            log.Printf("Get all funds: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        descption, err := GetDescription(fundName)
        if err != nil {
            log.Printf("Get all funds: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        err = bondListTmpl.Execute(w, FundDetail{fundName, descption, bonds})
        //fmt.Printf("funds: %#v\n", funds)
        if err != nil {
            log.Printf("Transefor to bondListTmpl html error: %v", err)
        }
    case "PUT":
        bond := r.FormValue("bond")
        //fmt.Printf("bondsHandler insert %s to %s \n", bond, fundName)
        if err:= InsertBond(fundName, bond); err != nil{
            log.Printf("Insert bond error: %v", err)
            http.Error(w, "Insert bond error", http.StatusInternalServerError)
        }else{
	    http.Error(w, "bond inserted Successfully!", http.StatusOK)
	}
    case "POST":
	err := r.ParseForm()
	if err != nil {
            fmt.Println("err ", err)
	    http.Error(w, "parse failed", http.StatusInternalServerError)
	}
	bond := r.Form["bond"][0]
        //fmt.Printf("bondsHandler delete %s from %s \n", bond, fundName)
        if err:= RemoveBond(fundName, bond); err != nil{
            log.Printf("Remove bond error: %v", err)
            http.Error(w, "Remove bond error", http.StatusInternalServerError)
        } else{
	    http.Error(w, "Bond reomved successfully!", http.StatusOK)
	}
    default:
        http.Error(w, fmt.Sprintf("HTTP Method %s Not Allowed", r.Method), http.StatusMethodNotAllowed)
    }
}
func getAllFundHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("getAllFundHandler called method %v \n", r.Method)
    switch r.Method {
    case "GET":
        funds, err := GetAllfunds()
        if err != nil {
            log.Printf("Get all funds: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        err = indexTmpl.Execute(w, funds)
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
