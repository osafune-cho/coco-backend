package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn, err := getDsn()
	if err != nil {
		return err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.DB()
	db.AutoMigrate(&Team{})
	DB = db
	return nil
}

func getEnvOrError(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing environment variable %s", key)
	}
	return value, nil
}

func getDsn() (string, error) {
	host, err := getEnvOrError("DB_HOST")
	if err != nil {
		return "", err
	}

	port, err := getEnvOrError("DB_PORT")
	if err != nil {
		return "", err
	}

	user, err := getEnvOrError("DB_USER")
	if err != nil {
		return "", err
	}

	password, err := getEnvOrError("DB_PASSWORD")
	if err != nil {
		return "", err
	}

	dbName, err := getEnvOrError("DB_NAME")
	if err != nil {
		return "", err
	}

	config := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host,
		user,
		password,
		dbName,
		port,
		"disable",
		"Asia/Tokyo",
	)

	return config, nil
}
