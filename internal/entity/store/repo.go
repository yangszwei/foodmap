package store

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IStoreRepo store repository interface
type IStoreRepo interface {
	InsertOne(store Store) (primitive.ObjectID, error)
	FindOneByID(id primitive.ObjectID, fields bson.M) (Store, error)
	Find(query string, categories []string, fields bson.M, limit, skip int64) ([]Store, error)
	UpdateOne(store Store) error
	DeleteOne(id primitive.ObjectID) error
	InsertComment(id primitive.ObjectID, comment Comment) (primitive.ObjectID, error)
	FindComments(storeID primitive.ObjectID, limit, skip int64) ([]Comment, error)
	DeleteComment(storeID, commentID primitive.ObjectID) error
}
