package server

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"fmt"
)

func Db_connect(db Db_data) *sql.DB {
	dsn := build_dsn(db)
	dbh, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil
	}
	err2 := dbh.Ping()
	if err2 != nil {
		return nil
	}
	return dbh
}

func Db_disconnect(dbh *sql.DB) {
	dbh.Close()
}

func build_dsn(db Db_data) string {
	address := fmt.Sprintf("%v:%d", db.Host, db.Port)
	config := mysql.Config{
		User: db.User,
		Passwd: db.Password,
		Net: "tcp",
		Addr: address,
		DBName: db.Database,
	}
	return config.FormatDSN()
}

