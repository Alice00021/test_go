package models

import (
    "fmt"
    "test_go/config"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func Init_DB(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimeZone,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
    }

    return db, nil
}