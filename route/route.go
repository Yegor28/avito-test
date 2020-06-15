package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Yegor28/avito-test/entity"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strings"
	//"os"
)

	//var host    = os.Getenv("DB_HOST")
	//var port     = os.Getenv("DB_PORT")
	//var user     = os.Getenv("DB_USER")
	//var password = os.Getenv("DB_PASSWORD")
	//var dbname   = os.Getenv("DB_NAME")


var host    = "http://localhost"
var port     = "5432"
var user     = "postgres"
var password = "452814"
var dbname   = "ads"


type errInfo struct {
	Err  error  `json:"error"`
	Info string `json:"info"`
}

func dbConnect(host, port, user, password, dbname string) (*sql.DB, error) {
	//psqlInfo := fmt.Sprintf("user=%s "+"password=%s dbname=%s port=%s sslmode=disable", user, password, dbname, port)
	psqlInfo := "host=db port=5432 user=postgres password=452814 dbname=ads sslmode=disable"
	db, _ := sql.Open("postgres", psqlInfo)

	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}

func CreateAd(w http.ResponseWriter, r *http.Request) {
	var advert entity.Ad
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
		fmt.Printf("Can't connect to database. Error in function dbConnect. Param's function host: %s,"+
			" user: %s, port: %s, password: %s, dbname: %s \n", host, user, port, password, dbname)
		info.Err = err
		info.Info = "Can't connect to database. Invalid data."
		json.NewEncoder(w).Encode(info)
		return
	}
	defer db.Close()

	statement := "INSERT INTO advert(name, description, photos, price) " +
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
func GetAd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db, err := dbConnect(host, port, user, password, dbname)
	if err != nil {
		var info errInfo
		fmt.Printf("Can't connect to database. Error in function dbConnect. Param's function host: %s,"+
			" user: %s, port: %s, password: %s, dbname: %s \n", host, user, port, password, dbname)
		info.Err = err
		info.Info = "Can't connect to database. Invalid data."
		json.NewEncoder(w).Encode(info)
		return
	}
	defer db.Close()

	statement := "SELECT id, name, description, photos, price FROM advert where id=$1;"
	stmt, err := db.Prepare(statement)
	if err != nil {
		var info errInfo
		fmt.Printf("Can't prepare query. Error in db.Prepare function. Param's function %s.\n", statement)
		info.Err = err
		info.Info = "Can't prepare query"
		json.NewEncoder(w).Encode(info)
		return
	}
	var advert entity.Ad
	stmt.QueryRow(vars["id"]).Scan(&advert.Id, &advert.Name, &advert.Description, pq.Array(&advert.Photos), &advert.Price)
	if advert.Id == 0 && advert.Name == "" && advert.Description == "" && len(advert.Photos) == 0 && advert.Price == 0 {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `"error":"ad does not exist"`)
		fmt.Println("Ad does not exist")
		return
	}
	fields := r.URL.Query().Get("fields")
	switch fields {
	case "description":
		w.Header().Set("Content-Type", "application/json")
		if (len(advert.Photos) == 0) {
			fmt.Fprintf(w, `{"Name": "%s", "Description":"%s", "Price": %d, "Photo": "%s"}`,
				advert.Name, advert.Description, advert.Price, advert.Photos)
			w.WriteHeader(http.StatusOK)
			return
		}
		fmt.Fprintf(w, `{"Name": "%s", "Description":"%s", "Price": %d, "Photo": "%s"}`,
			advert.Name, advert.Description, advert.Price, advert.Photos[0])
		return
	case "photos":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"Name": "%s", "Price": %d, "Photos": [%+q]}`,
			advert.Name, advert.Price, strings.Join(advert.Photos, ", "))
		return
	case "all":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(advert)
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		if (len(advert.Photos) == 0) {
			fmt.Fprintf(w, `{"Name": "%s", "Price": %d, "Photo": "%s"}`, advert.Name, advert.Price, advert.Photos)
			w.WriteHeader(http.StatusOK)
			return
		}
		fmt.Fprintf(w, `{"Name": "%s", "Price": %d, "Photo": "%s"}`, advert.Name, advert.Price, advert.Photos[0])
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Println("YES")
}
func GetAdsList(w http.ResponseWriter, r *http.Request) {
	var params entity.SearchRequest
	params.Limit = 10
	var allAds []entity.Ad
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &params)
	if err != nil {
		fmt.Printf("Hi")
		fmt.Printf("Invalid params")
		return
	}

	db, err := dbConnect(host, port, user, password, dbname)
	if err != nil {
		var info errInfo
		fmt.Printf("Can't connect to database. Error in function dbConnect. Param's function host: %s,"+
			" user: %s, port: %s, password: %s, dbname: %s \n", host, user, port, password, dbname)
		info.Err = err
		info.Info = "Can't connect to database. Invalid data."
		json.NewEncoder(w).Encode(info)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM advert;")
	for rows.Next() {
		var row entity.Ad
		rows.Scan(&row.Id, &row.Name, &row.Description, pq.Array(&row.Photos), &row.Price, &row.Time)
		allAds = append(allAds, row)
	}
	//сортировка
	if params.OrderBy != 0 {
		switch params.OrderField {
		case "price":
			sort.SliceStable(allAds, func(i, j int) bool {
				return allAds[i].Price < allAds[j].Price && (params.OrderBy == 1)
			})
		case "time":
			sort.SliceStable(allAds, func(i, j int) bool {
				return allAds[i].Time.After(allAds[j].Time) && (params.OrderBy == 1)
			})
		}
	}

	//Пагинация
	var pageNum float64
	var i int
	pageNum = math.Ceil(float64(len(allAds)) / float64(params.Limit))
	var pagesArr []entity.Page
	fmt.Println(pageNum)
	for i = 0; i < int(pageNum); i++ {
		start := int(math.Min(float64(i*params.Limit), float64(len(allAds))))
		end := int(math.Min(float64((i+1)*params.Limit), float64(len(allAds))))
		page := entity.Page{Page_number: i + 1, Page_size: len(allAds[start:end]), Adverts: allAds[start:end]}
		pagesArr = append(pagesArr, page)
	}

	if params.Page != 0 {
		resp, err := json.Marshal(pagesArr[params.Page-1])
		if err != nil {
			var info errInfo
			fmt.Printf("Can't marshal pages\n")
			info.Err = err
			info.Info = "Can't marshal pages"
			json.NewEncoder(w).Encode(info)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(resp))
		return
	}

	resp, err := json.Marshal(pagesArr)
	if err != nil {
		var info errInfo
		fmt.Printf("Can't marshal pagesArr\n")
		info.Err = err
		info.Info = "Can't marshal pages"
		json.NewEncoder(w).Encode(info)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(resp))

}
