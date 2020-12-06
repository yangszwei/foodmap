package persistence_test

import (
	"foodmap/internal/infra/config"
	"foodmap/internal/infra/persistence"
	"testing"
)

func Open() (persistence.DB, error) {
	cfg, err := config.Open("../../../.env")
	if err != nil {
		return persistence.DB{}, err
	}
	db, err := persistence.Open(cfg.DB)
	if err != nil {
		return persistence.DB{}, err
	}
	return db, nil
}

func TestDB(t *testing.T) {
	cfg, err := config.Open("../../../.env")
	if err != nil {
		t.Error(err)
	}
	db, err := persistence.Open(cfg.DB)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Error(err)
	}
}
