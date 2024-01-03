package db

import (
	"os"

	"github.com/Stolkerve/kappa/pgk/storage"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() {
	env := os.Getenv("ENV")
	if env == "DEV" {
		if db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		}); err != nil {
			panic(err)
		} else {
			DB = db
		}
	} else {
		dsn := os.Getenv("DSN")
		if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		}); err != nil {
			panic(err)
		} else {
			DB = db
		}
	}
	DB.AutoMigrate(&storage.Function{}, &storage.Call{})
}
