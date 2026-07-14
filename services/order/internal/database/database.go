package database

import (
	"fmt"
	"log"
    "github.com/rudransh/distributed-commerce/order/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	cfg := DBConfig

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&model.Order{},
		&model.OrderItem{},
		&model.IdempotencyKey{},
		&model.OutboxEvent{},
	)

	if err != nil {
		log.Fatal(err)
	}

	DB = db

	log.Println("Database Connected")

}
