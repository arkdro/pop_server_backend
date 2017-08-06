package server

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func Run(web_port int, db Db_data) {
	dbh := Db_connect(db)
	handler := Handler{dbh}
	http.Handle("/point", handler)
	http.HandleFunc("/countries",  get_countries)
	address := ":" + strconv.Itoa(web_port)
	res := http.ListenAndServe(address, nil)
	log.Fatal(res)
	Db_disconnect(dbh)
}

type Handler struct {
	db *sql.DB
}

func (h Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	countries := get_countries_from_database(h.db)
	res := build_countries_response(countries)
	writer.Write(res)
}

func get_countries(writer http.ResponseWriter, request *http.Request) {
}

func get_countries_from_database(db *sql.DB) []string {
	countries := make([]string, 0)
	rows, err := db.Query("select country from countries")
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

func build_countries_response(countries []string) []byte {
	return nil
}

