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
		EmailConfig      EmailConfig
		JWT              JWTConfig
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
		BasePath   string `env:"LOCAL_FILE_STORAGE_BASE_PATH,required"`
		ExportPath string `env:"LOCAL_FILE_STORAGE_EXPORT_PATH,required"`
	}
	EmailConfig struct {
		SMTPHost       string `env:"SMTP_HOST,required"`
		SMTPPort       int    `env:"SMTP_PORT,required"`
		SenderEmail    string `env:"SENDER_EMAIL,required"`
		SenderPassword string `env:"SENDER_PASSWORD,required"`
		VerifyBaseURL  string `env:"VERIFY_BASE_URL,required"`
	}
	JWTConfig struct {
		SecretKey string `env:"JWT_SECRET_KEY,required"`
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
