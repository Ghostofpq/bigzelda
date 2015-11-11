package main

import (
"fmt"
"net/http"
"github.com/gorilla/mux"
)

func main() {
    fmt.Println("hello boyyy")
	r := mux.NewRouter()
	r.HandleFunc("/gorilla", GorillaHandler)
    //r.HandleFunc("/{value}", RedirectionHandler)
    //r.HandleFunc("/shortlink/{value}", ShortlinkCreationHandler)
    //r.HandleFunc("/admin/{value}", MonitoringHandler)
    http.Handle("/", r)

	http.ListenAndServe(":8000", r)
}

func GorillaHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Gorilla!\n"))
}