package store

import (
	"foodmap/internal/entity/store"
	"foodmap/internal/infra/object"
	"foodmap/internal/infra/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// NewRepo create an instance of Repo
func NewRepo(db *persistence.DB) (*Repo, error) {
	s := new(Repo)
	s.db = db
	s.stores = db.DB.Collection("stores")
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
	_, err := s.stores.Indexes().CreateOne(s.db.Ctx, index)
	return s, err
}

// Repo implement store.IStoreRepo
type Repo struct {
	db *persistence.DB
	stores *mongo.Collection
}

// InsertOne insert a store record to database
func (r Repo) InsertOne(store store.Store) (primitive.ObjectID, error) {
	store.ID = primitive.NewObjectID()
	t := time.Now()
	store.CreatedAt = &t
	store.UpdatedAt = &t
	_, err := r.stores.InsertOne(r.db.Ctx, store)
	return store.ID, err
}

// FindOneByID find a store record by ID
func (r Repo) FindOneByID(id primitive.ObjectID, fields bson.M) (store.Store, error) {
	var result store.Store
	opt := options.FindOne()
	if fields != nil {
		opt = opt.SetProjection(fields)
	}
	err := r.stores.FindOne(r.db.Ctx, bson.M{ "_id": id }, opt).Decode(&result)
	return result, err
}

// Find find a list of store records from database
func (r Repo) Find(query string, categories []string, fields bson.M, limit, skip int64) ([]store.Store, error) {
	var ( f = make(bson.M) ; result []store.Store )
	if query != "" {
		f = bson.M{ "$text": bson.M{ "$search": query } }
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
	cursor, err := r.stores.Find(r.db.Ctx, f, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.db.Ctx)
	for cursor.Next(r.db.Ctx) {
		var record store.Store
		if err := cursor.Decode(&record) ; err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// UpdateOne update a store repository on database
func (r Repo) UpdateOne(store store.Store) error {
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
	_, err = r.stores.UpdateOne(
		r.db.Ctx,
		bson.M{ "_id": store.ID },
		bson.M{ "$set": update },
	)
	return err
}

// DeleteOne delete a store record from database
func (r Repo) DeleteOne(id primitive.ObjectID) error {
	_, err := r.stores.DeleteOne(r.db.Ctx, bson.M{ "_id": id })
	return err
}

// InsertComment add a comment to store record
func (r Repo) InsertComment(id primitive.ObjectID, comment store.Comment) (primitive.ObjectID, error) {
	var update bson.M
	comment.ID = primitive.NewObjectID()
	t := time.Now()
	comment.CreatedAt = &t
	data, err := bson.Marshal(comment)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	if err = bson.Unmarshal(data, &update) ; err != nil {
		return primitive.ObjectID{}, err
	}
	_, err = r.stores.UpdateOne(
		r.db.Ctx,
		bson.M{ "_id": id },
		bson.M{ "$push": bson.M{ "cmnt": update } },
	)
	return comment.ID, err
}

// FindComments find a list of comments from store record
func (r Repo) FindComments(storeID primitive.ObjectID, limit, skip int64) ([]store.Comment, error) {
	var result store.Store
	opt := options.FindOne().SetProjection(bson.M{
		"cmnt": bson.M{
			"$slice": []int64{ skip, limit },
		},
	}).SetProjection(bson.M{
		"cmnt": true,
	})
	err := r.stores.FindOne(r.db.Ctx, bson.M{
		"_id": storeID,
	}, opt).Decode(&result)
	log.Println(result)
	return result.Comments, err
}

// DeleteComment remove a comment by ID from store record
func (r Repo) DeleteComment(storeID, commentID primitive.ObjectID) error {
	_, err := r.stores.UpdateOne(r.db.Ctx, bson.M{ "_id": storeID }, bson.M{
		"$pull": bson.M{
			"cmnt": bson.M{ "_id": commentID },
		},
	})
	return err
}
