package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // here
)

var db *sql.DB

var (
	host     = "localhost"
	port     = 25434
	user     = "docker"
	password = "docker"
	dbname   = "gis"
)

func Connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	dbConnection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	db = dbConnection
}

func GetDb() *sql.DB {
	return db
}
