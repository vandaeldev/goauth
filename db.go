package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=root dbname=goauth port=5432 sslmode=disable TimeZone=Europe/Amsterdam"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection failed: %s", err.Error())
	}
	db = conn
	migrateDB()
}

func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close database: %s", err.Error())
	} else {
		sqlDB.Close()
	}
}

func migrateDB() {
	db.AutoMigrate(&User{})
}
