package config

import "fmt"

// DBConfig database config
type DBConfig struct {
	Name     string `json:"DB_NAME" validate:"required"`
	Host     string `json:"DB_HOST" validate:"required"`
	Port     string `json:"DB_PORT" validate:"required"`
	User     string `json:"DB_USER" validate:"required"`
	Password string `json:"DB_PASSWORD" validate:"required"`
}

// ToURI convert DBConfig to URI string
func (d *DBConfig) ToURI() string {
	f := "mongodb://%s:%s@%s:%s/%s?authSource=admin"
	return fmt.Sprintf(f, d.User, d.Password, d.Host, d.Port, d.Name)
}
