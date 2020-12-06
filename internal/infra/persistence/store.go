package persistence

import (
	"foodmap/internal/entity/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// NewStoreRepo create an instance of StoreRepo
func NewStoreRepo(db *DB) (*StoreRepo, error) {
	s := new(StoreRepo)
	s.db = db
	s.stores = db.db.Collection("stores")
	index := mongo.IndexModel{ Keys:
		bson.D{
			bson.E{ Key: "name", Value: "text" },
			bson.E{ Key: "desc", Value: "text" },
			bson.E{ Key: "cat", Value: "text" },
			bson.E{ Key: "menu.name", Value: "text" },
			bson.E{ Key: "menu.desc", Value: "text" },
			bson.E{ Key: "menu.cat", Value: "text" },
		},
	}
	_, err := s.stores.Indexes().CreateOne(s.db.ctx, index)
	return s, err
}

// StoreRepo implement store.IStoreRepo
type StoreRepo struct {
	db *DB
	stores *mongo.Collection
}

// InsertOne insert a store record to database
func (s StoreRepo) InsertOne(store store.Store) (primitive.ObjectID, error) {
	store.ID = primitive.NewObjectID()
	t := time.Now()
	store.CreatedAt = &t
	store.UpdatedAt = &t
	_, err := s.stores.InsertOne(s.db.ctx, store)
	return store.ID, err
}

// FindOneByID find a store record by ID
func (s StoreRepo) FindOneByID(id primitive.ObjectID, fields bson.M) (store.Store, error) {
	var result store.Store
	opt := options.FindOne()
	if fields != nil {
		opt = opt.SetProjection(fields)
	}
	err := s.stores.FindOne(s.db.ctx, bson.M{ "_id": id }, opt).Decode(&result)
	return result, err
}

// Find find a list of store records from database
func (s StoreRepo) Find(query string, filter store.Store, fields bson.M, limit, skip int64) ([]store.Store, error) {
	var ( f bson.M ; result []store.Store )
	data, err := bson.Marshal(filter)
	if err != nil {
		return nil, err
	}
	if query != "" {
		f = bson.M{ "$text": bson.M{ "$search": query } }
	}
	if err = bson.Unmarshal(data, &f) ; err != nil {
		return nil, err
	}
	opt := options.Find()
	opt.SetLimit(limit)
	opt.SetSkip(skip)
	if fields != nil {
		opt = opt.SetProjection(fields)
	}
	cursor, err := s.stores.Find(s.db.ctx, f, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.db.ctx)
	for cursor.Next(s.db.ctx) {
		var record store.Store
		if err := cursor.Decode(&record) ; err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// UpdateOne update a store repository on database
func (s StoreRepo) UpdateOne(store store.Store) error {
	var update bson.D
	t := time.Now()
	store.UpdatedAt = &t
	data, err := bson.Marshal(store)
	if err != nil {
		return err
	}
	if err := bson.Unmarshal(data, &update) ; err != nil {
		return err
	}
	_, err = s.stores.UpdateOne(
		s.db.ctx,
		bson.M{ "_id": store.ID },
		bson.M{ "$set": update },
	)
	return err
}

// DeleteOne delete a store record from database
func (s StoreRepo) DeleteOne(id primitive.ObjectID) error {
	_, err := s.stores.DeleteOne(s.db.ctx, bson.M{ "_id": id })
	return err
}

// InsertComment add a comment to store record
func (s StoreRepo) InsertComment(id primitive.ObjectID, comment store.Comment) error {
	var update bson.M
	comment.ID = primitive.NewObjectID()
	t := time.Now()
	comment.CreatedAt = &t
	data, err := bson.Marshal(comment)
	if err != nil {
		return err
	}
	if err = bson.Unmarshal(data, &update) ; err != nil {
		return err
	}
	_, err = s.stores.UpdateOne(
		s.db.ctx,
		bson.M{ "_id": id },
		bson.M{ "$push": bson.M{ "cmnt": update } },
	)
	return err
}

// FindComments find a list of comments from store record
func (s StoreRepo) FindComments(storeID primitive.ObjectID, limit, skip int64) ([]store.Comment, error) {
	var result store.Store
	opt := options.FindOne().SetProjection(bson.M{
		"cmnt": bson.M{
			"$slice": []int64{ skip, limit },
		},
	}).SetProjection(bson.M{
		"cmnt": true,
	})
	err := s.stores.FindOne(s.db.ctx, bson.M{
		"_id": storeID,
	}, opt).Decode(&result)
	log.Println(result)
	return result.Comments, err
}

// DeleteComment remove a comment by ID from store record
func (s StoreRepo) DeleteComment(storeID, commentID primitive.ObjectID) error {
	_, err := s.stores.UpdateOne(s.db.ctx, bson.M{ "_id": storeID }, bson.M{
		"$pull": bson.M{
			"cmnt": bson.M{ "_id": commentID },
		},
	})
	return err
}
