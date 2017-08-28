package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
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
	write_response_countries(writer, req, res)
}

func write_response_countries(writer http.ResponseWriter, req *http.Request, data Countries) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func write_response_point(writer http.ResponseWriter, req *http.Request, data Point) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func write_response_points(writer http.ResponseWriter, req *http.Request, data Points) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		log.Printf("json encode points error: %v", err)
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
	country := vars["country"]
	log.Printf("get_country, %v", country)
	data := build_data_points(dbh, country)
	write_response_points(writer, request, data)
}

func get_point(writer http.ResponseWriter, request *http.Request, dbh *sql.DB) {
	log.Printf("get_point")
	vars := mux.Vars(request)
	country := vars["country"]
	year_str := vars["year"]
	log.Printf("get_point, %v, %v", country, year_str)
	year, err := strconv.Atoi(year_str)
	var data Point
	if err != nil {
		log.Printf("get_point, wrong year, %v", year_str)
	} else {
		data = build_data_point(dbh, country, year)
	}
	write_response_point(writer, request, data)
}

func build_data_points(db *sql.DB, country string) Points {
	var points Points
	cmd := "SELECT country_code, year, age FROM country_median_age where country = ?"
	rows, err := db.Query(cmd, country)
	if err != nil {
		return points
	}
	points = extract_data_points_from_db_result(rows)
	return points
}

func build_data_point(db *sql.DB, country string, year int) Point {
	var data Point
	date := build_date(year)
	cmd := "SELECT country_code, age FROM country_median_age where country = ? and year = ?"
	row := db.QueryRow(cmd, country, date)
	data = extract_data_point_from_db_result(year, row)
	return data
}

func build_date(year int) string {
	return strconv.Itoa(year) + "-00-00"
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

func extract_data_points_from_db_result(rows *sql.Rows) Points {
	points := make(Points, 0)
	for rows.Next() {
		var code string
		var date string
		var age float64
		err := rows.Scan(&code, &date, &age)
		if err != nil {
			log.Printf("extract one of points error: %v", err)
			continue
		}
		point := build_one_point(date, age)
		points = append(points, point)
	}
	return points
}

func build_one_point(date string, age float64) Point {
	log.Printf("build_one_point: %v, %v", date, age)
	re := regexp.MustCompile("^(\\d+)-")
	years := re.FindStringSubmatch(date)
	year, err := strconv.Atoi(years[1])
	var point Point
	if err != nil {
		log.Printf("build_one_point, wrong date, %v", years)
	} else {
		point = Point{
			Year: year,
			Value: age,
		}
	}
	return point
}

func extract_data_point_from_db_result(year int, row *sql.Row) Point {
	var data Point
	var code string
	var age float64
	err := row.Scan(&code, &age)
	if err != nil {
		log.Printf("extract data point error: %v", err)
	} else {
		data = Point{
			Year: year,
			Value: age,
		}
	}
	return data
}

func build_countries_response(countries []string) Countries {
	res := make(Countries, 0)
	for _, country := range countries {
		res = append(res, Country{country})
	}
	return res
}

