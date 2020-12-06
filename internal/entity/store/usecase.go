package store

import "foodmap/internal/infra/object"

// IStoreUsecase store usecase interface
type IStoreUsecase interface {
	CreateOne(document object.H) (string, error)
	FindOneByID(id string, fields string) (object.H, error)
	Find(query string, filter object.H, fields string, limit, skip int) ([]object.H, error)
	UpdateOne(document object.H) error
	DeleteOne(id string) error
	CreateComment(storeID string, comment object.H) error
	FindComments(storeID string, limit, skip int) ([]object.H, error)
	DeleteComment(storeID, commentID string) error
}
