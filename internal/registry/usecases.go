package registry

import (
	"foodmap/internal/entity/store"
	storeInfra "foodmap/internal/infra/entity/store"
)

type usecases struct {
	store store.IStoreUsecase
}

func setupUsecase(r repos) (u usecases) {
	u.store = storeInfra.NewUsecase(r.store, v)
	return
}
