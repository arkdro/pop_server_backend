package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Run(web_port int, db Db_data) {
	dbh := Db_connect(db)
	r := mux.NewRouter()
	r.HandleFunc("/point/{country}/{year}",
		func(writer http.ResponseWriter, req *http.Request) {
			get_point(writer, req, dbh)
		})
	r.HandleFunc("/country",
		func(writer http.ResponseWriter, req *http.Request) {
			get_countries(writer, req, dbh)
		})
	r.HandleFunc("/country/{country}",
		func(writer http.ResponseWriter, req *http.Request) {
			get_country(writer, req, dbh)
		})
	address := ":" + strconv.Itoa(web_port)
	res := http.ListenAndServe(address, r)
	log.Fatal(res)
	Db_disconnect(dbh)
}

type Point_handler struct {
	db *sql.DB
}

func (h Point_handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	countries := get_countries_from_database(h.db)
	res := build_countries_response(countries)
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(res)
	if err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func get_countries(writer http.ResponseWriter, request *http.Request, dbh *sql.DB) {
	log.Printf("get_countries")
	handler := Point_handler{dbh}
	handler.ServeHTTP(writer, request)
}

func get_country(writer http.ResponseWriter, request *http.Request, dbh *sql.DB) {
	log.Printf("get_country")
	vars := mux.Vars(request)
	c := vars["country"]
	log.Printf("get_country, %v", c)
	handler := Point_handler{dbh}
	handler.ServeHTTP(writer, request)
}

func get_point(writer http.ResponseWriter, request *http.Request, dbh *sql.DB) {
	log.Printf("get_point")
	vars := mux.Vars(request)
	country := vars["country"]
	year := vars["year"]
	log.Printf("get_point, %v, %v", country, year)
	handler := Point_handler{dbh}
	handler.ServeHTTP(writer, request)
}

func get_countries_from_database(db *sql.DB) []string {
	countries := make([]string, 0)
	cmd := "SELECT DISTINCT country FROM country_median_age ORDER BY country"
	rows, err := db.Query(cmd)
	if err != nil {
		return countries
	}
	defer rows.Close()
	countries = extract_countries_from_db_result(rows)
	err = rows.Err()
	if err != nil {
		log.Printf("get countries db error: %v", err)
	}
	return countries
}

func extract_countries_from_db_result(rows *sql.Rows) []string {
	countries := make([]string, 0)
	for rows.Next() {
		var country string
		err := rows.Scan(&country)
		if err != nil {
			log.Printf("extract country error: %v", err)
			continue
		}
		countries = append(countries, country)
	}
	return countries
}

func build_countries_response(countries []string) Countries {
	res := make(Countries, 0)
	for _, country := range countries {
		res = append(res, Country{country})
	}
	return res
}

