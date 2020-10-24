package main

import (
	"fmt"
	"github.com/Yegor28/avito-test/route"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

var host    = "db"
var port     = "5432"
var user     = "postgres"
var password = "452814"
var dbname   = "ads"

func main() {
	r := mux.NewRouter()
	db, _:= route.DbConnect(host, port, user, password, dbname)

	fmt.Printf("Table was created")
	defer db.Close()

	statement := "CREATE TABLE IF NOT EXISTS advert (id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY, name varchar(255), description varchar(255), photos text[], price integer);"
	stmt, _ := db.Prepare(statement)
	defer stmt.Close()
	stmt.QueryRow()
	db.Query("")

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
