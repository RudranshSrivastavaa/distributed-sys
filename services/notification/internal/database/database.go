package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {

	dsn := "host=localhost user=postgres password=postgres dbname=notificationdb port=5435 sslmode=disable"

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		return err
	}

	DB = db

	log.Println("Database Connected")
	return nil

}