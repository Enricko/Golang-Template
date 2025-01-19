package config

import (
    "fmt"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

func InitDB(config *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        config.DBUser,
        config.DBPassword,
        config.DBHost,
        config.DBPort,
        config.DBName,
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    return db, nil
}