package config

import (
	"encoding/json"
	"foodmap/internal/infra/validator"
	"github.com/joho/godotenv"
)

// Open load configuration from path
func Open(path string) (cfg Config, err error) {
	var raw map[string]string
	raw, err = godotenv.Read(path)
	if err != nil {
		return
	}
	ser, _ := json.Marshal(raw)
	_ = json.Unmarshal(ser, &cfg.DB)
	_ = json.Unmarshal(ser, &cfg.Server)
	return
}

// Config configuration
type Config struct {
	DB     DBConfig
	Server ServerConfig
}

// Validate validate configuration
func Validate(c Config, v validator.Validator) error {
	return v.Validate(c)
}
