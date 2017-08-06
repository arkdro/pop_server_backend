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
	lst := get_countries_from_database(h.db)
	writer.Write(lst)
}

func get_countries(writer http.ResponseWriter, request *http.Request) {
}

func get_countries_from_database(db *sql.DB) []byte {
	res := make([]byte, 0)
	res = append(res, []byte("c1\n")...)
	res = append(res, []byte("c2\n")...)
	return res
}

