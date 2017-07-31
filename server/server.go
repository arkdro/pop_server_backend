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
}

func (h Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
}

func get_countries(resp http.ResponseWriter, req *http.Request) {
}

