package server

import (
	"log"
	"net/http"
)

func Run(db Db_data) {
	handler := Handler{}
	http.Handle("/point", handler)
	http.HandleFunc("/countries",  get_countries)
	res := http.ListenAndServe(":8080", nil)
	log.Fatal(res)
}

type Handler struct {
}

func (h Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
}

func get_countries(resp http.ResponseWriter, req *http.Request) {
}

