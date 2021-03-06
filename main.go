package main

import (
	"github.com/asdf/pop_server/server"
	"flag"
)

func main() {
	var web_port = flag.Int("web_port", 8080, "web port")
	var host = flag.String("host", "localhost", "db host")
	var port = flag.Int("port", 3306, "db port")
	var user = flag.String("user", "root", "db user")
	var password = flag.String("password", "root", "db password")
	var database = flag.String("database", "pdata", "database")
	flag.Parse()
	db_data := server.Db_data{
		Host: *host,
		Port: *port,
		User: *user,
		Password: *password,
		Database: *database,
	}
	server.Run(*web_port, db_data)
}

