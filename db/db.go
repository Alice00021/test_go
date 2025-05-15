package db

import (
    "fmt"
    "test_go/config"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func Init_DB(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
        cfg.Database.Host,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Name,
        cfg.Database.Port,
        cfg.Database.SSLMode,
        cfg.Database.TimeZone,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
    }

    return db, nil
}
