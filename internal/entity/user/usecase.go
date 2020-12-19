package user

import "foodmap/internal/infra/object"

// IStoreUsecase store usecase interface
type IUserUsecase interface {
	CreateOne(document object.H) (string, error)
	FindOneByID(id string, fields string) (object.H, error)
	Find(query, fields string, limit, skip int64) ([]object.H, error)
	UpdateOne(document object.H) error
	DeleteOne(id string) error
}
