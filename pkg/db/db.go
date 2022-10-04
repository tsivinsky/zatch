package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Connect() error {
	var err error

	dbHost := os.Getenv("DATABASE_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := "5432"

	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow", dbHost, dbUser, dbPassword, dbName, dbPort)

	Db, err = gorm.Open(postgres.Open(dbUri))

	err = Db.AutoMigrate(&Url{})

	return err
}
