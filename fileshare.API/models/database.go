package models

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var DB *sql.DB
var Q *Queries

func init() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("sqlite", dbUrl)
	if err != nil {
		panic(err)
	}
	DB = db
	Q = New(db)
}
