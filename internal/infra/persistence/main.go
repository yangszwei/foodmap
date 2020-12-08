/* Package persistence contains implementation of entity repositories */
package persistence

import (
	"context"
	"foodmap/internal/infra/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Open connect to database and return a DB instance
func Open(cfg config.DBConfig) (db DB, err error) {
	db.Ctx = context.Background()
	opt := options.Client().ApplyURI(cfg.ToURI())
	db.Client, err = mongo.Connect(db.Ctx, opt)
	if err != nil {
		return
	}
	db.DB = db.Client.Database(cfg.Name, nil)
	return
}

// DB wraps mongo client so the rest of the app don't need to know the
// implementation details
type DB struct {
	Client *mongo.Client
	DB     *mongo.Database
	Ctx    context.Context
}

// Close disconnect from database
func (db *DB) Close() error {
	return db.Client.Disconnect(db.Ctx)
}
