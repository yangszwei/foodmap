package registry

import (
	"foodmap/internal/entity/store"
	storeInfra "foodmap/internal/infra/entity/store"
	"foodmap/internal/infra/persistence"
)

type repos struct {
	store store.IStoreRepo
}

func setupRepos(db *persistence.DB) (r repos, err error) {
	r.store, err = storeInfra.NewRepo(db)
	if err != nil {
		return repos{}, err
	}
	return
}
