package usecase_test

import (
	"foodmap/internal/infra/config"
	"foodmap/internal/infra/persistence"
)

func Open() (persistence.DB, error) {
	cfg, err := config.Open("../../.env")
	if err != nil {
		return persistence.DB{}, err
	}
	db, err := persistence.Open(cfg.DB)
	if err != nil {
		return persistence.DB{}, err
	}
	return db, nil
}
