package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
	"github.com/go-playground/validator"

)

const (
	host = "127.0.0.1"
	port     = "5432"
	user     = "yegorp"
	password = "452814"
	dbname   = "ads"
)

type ad struct {
	Id int  `json: "ID"`
	Name        string `json:"Name" validate:"required,lte=200"`
	Description string `json:"Description" validate:"required,lte=1000"`
	Photos      []string `json:"Photos" validate:"required,lte=3"`
	Price       int `json:"Price" validate:"required"`
}

type errInfo struct {
	Err error `json:"error"`
	Info string `json:"info"`
}

func dbConnect(host, port, user, password, dbname string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, _ := sql.Open("postgres", psqlInfo)


	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}

func createAd(w http.ResponseWriter, r *http.Request)  {
	var advert ad
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &advert)
	if err != nil {
		fmt.Printf("Invalid params")
		return
	}

	validate := validator.New()
	err = validate.Struct(advert)
	if err != nil {
		var info errInfo
		fmt.Printf("Invalid params")
		info.Err = err
		info.Info = "Invalid params."
		json.NewEncoder(w).Encode(info)
		return
	}

	db, err := dbConnect(host, port, user, password, dbname)
	if err != nil {
		var info errInfo
		fmt.Printf("Can't connect to database. Error in function dbConnect. Param's function host: %s," +
			" user: %s, port: %s, password: s, dbname: %s \n", host, user, port, password, dbname)
		info.Err = err
		info.Info = "Can't connect to database. Invalid data."
		json.NewEncoder(w).Encode(info)
		return
	}
	defer db.Close()

	statement := "INSERT INTO advert(ad_name, description, photos, price) " +
		"VALUES ($1, $2, $3, $4);"
	stmt, err := db.Prepare(statement)
	if err != nil {
		var info errInfo
		fmt.Printf("Can't prepare query. Error in db.Prepare function. Param's function %s.\n", statement)
		info.Err = err
		info.Info = "Can't prepare query"
		json.NewEncoder(w).Encode(info)
		return
	}
	defer stmt.Close()

	stmt.QueryRow(advert.Name, advert.Description, pq.Array(advert.Photos), advert.Price)

	row, err := db.Query("SELECT id FROM advert ORDER BY id DESC LIMIT 1;")

	if err != nil {
		var info errInfo
		fmt.Printf("Can't make query. Error in db.Query function. Param's function %s.\n", "SELECT id FROM advert ORDER BY id DESC LIMIT 1;")
		info.Err = err
		info.Info = "Can't make query"
		json.NewEncoder(w).Encode(info)
		return
	}

	for row.Next() {
		row.Scan(&advert.Id)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"result":"success", "id": %d}`, advert.Id)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ads", createAd).Methods("POST")

	server := http.Server{
		Addr: ":228",
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("server run 127.0.0.1%s", server.Addr)
	server.ListenAndServe()
}
