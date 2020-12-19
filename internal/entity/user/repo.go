package user

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IUserRepo user repository interface
type IUserRepo interface {
	InsertOne(user User) (primitive.ObjectID, error)
	FindOneByID(id primitive.ObjectID, fields bson.M) (User, error)
	Find(query string, fields bson.M, limit, skip int64) ([]User, error)
	UpdateOne(user User) error
	DeleteOne(id primitive.ObjectID) error
}
