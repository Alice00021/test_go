package config 
import (
    "fmt"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    DBPort     string
    DBSSLMode  string
    DBTimeZone string
}

func Load() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, fmt.Errorf("ошибка загрузки .env файла: %v", err)
    }

    cfg := &Config{
        DBHost:     os.Getenv("DB_HOST"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
        DBPort:     os.Getenv("DB_PORT"),
        DBSSLMode:  os.Getenv("DB_SSLMODE"),
        DBTimeZone: os.Getenv("DB_TIMEZONE"),
    }

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" || cfg.DBPort == "" {
        return nil, fmt.Errorf("отсутствуют обязательные переменные окружения: DB_HOST, DB_USER, DB_NAME, DB_PORT")
    }
	return cfg, nil
}