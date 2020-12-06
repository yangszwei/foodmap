package config_test

import (
	"foodmap/internal/infra/config"
	"foodmap/internal/infra/validator"
	"testing"
)

var c config.Config

func TestOpen(t *testing.T) {
	var err error
	c, err = config.Open("../../../.env")
	if err != nil {
		t.Error(err)
	}
}

func TestValidate(t *testing.T) {
	if err := config.Validate(c, *validator.New()); err != nil {
		t.Error(err)
	}
}
