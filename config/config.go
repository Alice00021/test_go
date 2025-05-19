package config 
import (
    "fmt"
    "os"
    "gopkg.in/yaml.v2"

)

type Config struct {
    Database DatabaseConfig `yaml:"database"`
    Jwt JWTConfig `yaml:"jwt"`
}
type JWTConfig struct{
    SecretKey string `yaml:"SecretKey"`
}
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Name     string `yaml:"name"`
    Port     string `yaml:"port"`
    SSLMode  string `yaml:"sslmode"`
    TimeZone string `yaml:"timezone"`
}

func Load() (*Config, error) {
    filePath := "config/config.yaml"

    yamlFile, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("не удалось прочитать config.yaml: %v", err)
    }

    var cfg Config
    err = yaml.Unmarshal(yamlFile, &cfg)
    if err != nil {
        return nil, fmt.Errorf("не удалось распарсить YAML: %v", err)
    }

    if cfg.Database.Host == "" || cfg.Database.User == "" || cfg.Database.Name == "" || cfg.Database.Port == "" {
        return nil, fmt.Errorf("в config.yaml отсутствуют обязательные параметры БД")
    }

    if cfg.Jwt.SecretKey== ""  {
        return nil, fmt.Errorf("в config.yaml отсутствует секретный ключ для jwt")
    }

    return &cfg, nil
}
