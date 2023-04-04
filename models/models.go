package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var db, err = gorm.Open(sqlite.Open("data/files.db"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&File{})

	Db = db
}
