package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)


func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ads", createAd).Methods("POST")
	r.HandleFunc("/ads/{id:[0-9]+}", getAd).Methods("GET")
	r.HandleFunc("/ads", getAdsList).Methods("GET")
	server := http.Server{
		Addr: ":9999",
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
