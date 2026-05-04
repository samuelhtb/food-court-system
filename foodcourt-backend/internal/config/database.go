package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using default environment variables")
    }

    dsn := os.Getenv("SUPABASE_DB_URL")
    if dsn == "" {
        log.Fatal("Error: SUPABASE_DB_URL is not set in the environment variables")
    }

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }

    log.Println("Successfully connected!")

    DB = database

    return database
}