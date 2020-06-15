package main

import (
	"fmt"
	"github.com/Yegor28/avito-test/route"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ads", route.CreateAd).Methods("POST")
	r.HandleFunc("/ads/{id:[0-9]+}", route.GetAd).Methods("GET")
	r.HandleFunc("/ads", route.GetAdsList).Methods("GET")
	server := http.Server{
		Addr:         ":9999",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Server started. Port %s", server.Addr)
	server.ListenAndServe()
}
