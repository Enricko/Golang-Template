package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DBConnection string
    DBHost       string
    DBPort       string
    DBName       string
    DBUser       string
    DBPassword   string
    AppPort      string
    JWTSecret    string
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, err
    }

    config := &Config{
        DBConnection: os.Getenv("DB_CONNECTION"),
        DBHost:       os.Getenv("DB_HOST"),
        DBPort:       os.Getenv("DB_PORT"),
        DBName:       os.Getenv("DB_DATABASE"),
        DBUser:       os.Getenv("DB_USERNAME"),
        DBPassword:   os.Getenv("DB_PASSWORD"),
        AppPort:      os.Getenv("APP_PORT"),
        JWTSecret:    os.Getenv("JWT_SECRET"),
    }

    return config, nil
}