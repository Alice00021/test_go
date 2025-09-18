//package config
//
//import (
//	"fmt"
//	"gopkg.in/yaml.v2"
//	"os"
//)
//
//type Config struct {
//	Database DatabaseConfig `yaml:"database"`
//	Jwt      JWTConfig      `yaml:"jwt"`
//	SMTP     SMTPConfig     `yaml:"smtp"`
//}
//type JWTConfig struct {
//	SecretKey string `yaml:"SecretKey"`
//}
//type DatabaseConfig struct {
//	Host     string `yaml:"host"`
//	User     string `yaml:"user"`
//	Password string `yaml:"password"`
//	Name     string `yaml:"name"`
//	Port     string `yaml:"port"`
//	SSLMode  string `yaml:"sslmode"`
//	TimeZone string `yaml:"timezone"`
//}
//
//type SMTPConfig struct {
//	Host     string `yaml:"host"`
//	Port     int    `yaml:"port"`
//	Email    string `yaml:"email"`
//	Password string `yaml:"password"`
//}
//
//func Load() (*Config, error) {
//	filePath := "config/config.yaml"
//
//	yamlFile, err := os.ReadFile(filePath)
//	if err != nil {
//		return nil, fmt.Errorf("не удалось прочитать config.yaml: %v", err)
//	}
//
//	var cfg Config
//	err = yaml.Unmarshal(yamlFile, &cfg)
//	if err != nil {
//		return nil, fmt.Errorf("не удалось распарсить YAML: %v", err)
//	}
//
//	if cfg.Database.Host == "" || cfg.Database.User == "" || cfg.Database.Name == "" || cfg.Database.Port == "" {
//		return nil, fmt.Errorf("в config.yaml отсутствуют обязательные параметры БД")
//	}
//
//	if cfg.Jwt.SecretKey == "" {
//		return nil, fmt.Errorf("в config.yaml отсутствует секретный ключ для jwt")
//	}
//
//	if cfg.SMTP.Host == "" || cfg.SMTP.Port == 0 || cfg.SMTP.Email == "" || cfg.SMTP.Password == "" {
//		return nil, fmt.Errorf("в config.yaml отсутствуют обязательные параметры SMTP")
//	}
//
//	return &cfg, nil
//}

package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
)

type (
	// Config -.
	Config struct {
		App              App
		HTTP             HTTP
		Log              Log
		PG               PG
		Metrics          Metrics
		Swagger          Swagger
		LocalFileStorage LocalFileStorage
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	LocalFileStorage struct {
		BasePath string `env:"LOCAL_FILE_STORAGE_BASE_PATH,required"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Could not loading .env file")
	}

	// Parse environment variables into structs
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
