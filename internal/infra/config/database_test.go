package config_test

import (
	"foodmap/internal/infra/config"
	"testing"
)

func TestDBConfig_ToURI(t *testing.T) {
	uri := (&config.DBConfig{
		Name:       "name",
		Host:       "127.0.0.1",
		Port:       "27017",
		User:       "user",
		Password:   "password",
	}).ToURI()
	if uri != "mongodb://user:password@127.0.0.1:27017/name?authSource=admin" {
		t.Error(uri)
	}
}
