package store

import "foodmap/internal/infra/object"

// IStoreUsecase store usecase interface
type IStoreUsecase interface {
	CreateOne(document object.H) (string, error)
	FindOneByID(id string, fields string) (object.H, error)
	Find(query, categories, fields string, limit, skip int64) ([]object.H, error)
	UpdateOne(document object.H) error
	DeleteOne(id string) error
	CreateComment(storeID string, comment object.H) (string, error)
	FindComments(storeID string, admin bool, limit, skip int64) ([]object.H, error)
	DeleteComment(storeID, commentID string) error
}
