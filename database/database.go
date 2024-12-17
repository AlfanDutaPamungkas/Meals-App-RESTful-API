package database

import (
	"log"
	"meals-app/helper"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatabaseInit() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	helper.PanicError(err)

	log.Println("Connected to database")

	return db
}
