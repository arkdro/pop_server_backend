package server

import (
	"log"
	"net/http"
	"strconv"
)

func Run(web_port int, db Db_data) {
	handler := Handler{}
	http.Handle("/point", handler)
	http.HandleFunc("/countries",  get_countries)
	address := ":" + strconv.Itoa(web_port)
	res := http.ListenAndServe(address, nil)
	log.Fatal(res)
}

type Handler struct {
	db Db_data
}

func (h Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	lst := get_countries_from_database(h.db)
	writer.Write(lst)
}

func get_countries(writer http.ResponseWriter, request *http.Request) {
}

func get_countries_from_database(db Db_data) []byte {
	res := make([]byte, 0)
	res = append(res, []byte("c1\n")...)
	res = append(res, []byte("c2\n")...)
	return res
}

