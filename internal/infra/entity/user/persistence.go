package user

import (
	"foodmap/internal/entity/user"
	"foodmap/internal/infra/object"
	"foodmap/internal/infra/persistence"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewRepo create an instance of Repo
func NewRepo(db *persistence.DB) (*Repo, error) {
	s := new(Repo)
	s.db = db
	s.users = db.DB.Collection("users")
	index := mongo.IndexModel{Keys: bson.D{
		bson.E{Key: "name", Value: "text"},
		bson.E{Key: "email", Value: "text"},
	},
	}
	_, err := s.users.Indexes().CreateOne(s.db.Ctx, index)
	return s, err
}

// Repo implement store.IStoreRepo
type Repo struct {
	db    *persistence.DB
	users *mongo.Collection
}

// InsertOne insert a store record to database
func (r Repo) InsertOne(user user.User) (primitive.ObjectID, error) {
	user.ID = primitive.NewObjectID()
	t := time.Now()
	user.CreatedAt = &t
	user.UpdatedAt = &t
	_, err := r.users.InsertOne(r.db.Ctx, user)
	return user.ID, err
}

// FindOneByID find a store record by ID
func (r Repo) FindOneByID(id primitive.ObjectID, fields bson.M) (user.User, error) {
	var result user.User
	opt := options.FindOne()
	if fields != nil {
		opt = opt.SetProjection(fields)
	}
	err := r.users.FindOne(r.db.Ctx, bson.M{"_id": id}, opt).Decode(&result)
	return result, err
}

// Find find a list of store records from database
func (r Repo) Find(query string, categories []string, fields bson.M, limit, skip int64) ([]user.User, error) {
	var (
		f      = make(bson.M)
		result []user.User
	)
	if query != "" {
		f = bson.M{"$text": bson.M{"$search": query}}
	}
	if len(categories) != 0 {
		f["cat"] = object.H{
			"$all": categories,
		}
	}
	opt := options.Find().SetLimit(limit).SetSkip(skip)
	if fields != nil {
		opt = opt.SetProjection(fields)
	}
	cursor, err := r.users.Find(r.db.Ctx, f, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.db.Ctx)
	for cursor.Next(r.db.Ctx) {
		var record user.User
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// UpdateOne update a store repository on database
func (r Repo) UpdateOne(user user.User) error {
	var update bson.D
	t := time.Now()
	user.UpdatedAt = &t
	data, err := bson.Marshal(user)
	if err != nil {
		return err
	}
	if err := bson.Unmarshal(data, &update); err != nil {
		return err
	}
	_, err = r.users.UpdateOne(
		r.db.Ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": update},
	)
	return err
}

// DeleteOne delete a store record from database
func (r Repo) DeleteOne(id primitive.ObjectID) error {
	_, err := r.users.DeleteOne(r.db.Ctx, bson.M{"_id": id})
	return err
}
