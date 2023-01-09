package app

import (
	"database/sql"

	"github.com/hashicorp/go-hclog"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	log := hclog.Default()

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Error("Couldn't connect sqlite database ", err)
		panic(err)
	}
	return db
}
