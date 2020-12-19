package user_test

import (
	"foodmap/internal/entity/user"
	"foodmap/internal/infra/config"
	infra "foodmap/internal/infra/entity/user"
	"foodmap/internal/infra/persistence"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	testUser = user.User{
		Name:  "Test User",
		Email: "admin@gmail.com",
	}
	testUserID   primitive.ObjectID
	testUserRepo *infra.Repo
)

func TestNewUserRepo(t *testing.T) {
	cfg, err := config.Open("../../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
	db, err := persistence.Open(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}
	testUserRepo, err = infra.NewRepo(&db)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestUserRepo_InsertOne(t *testing.T) {
	var err error
	testUserID, err = testUserRepo.InsertOne(testUser)
	if err != nil {
		t.Error(err)
	}
}

func TestUserRepo_FindOneByID(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		s, err := testUserRepo.FindOneByID(testUserID, bson.M{"name": 1})
		if err != nil {
			t.Error(err)
		}
		if s.Name != testUser.Name {
			t.Error(s)
		}
		t.Log(s)
	})
}

func TestUserRepo_Find(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		query := []string{
			"store",
			"food",
			"Product",
			"Main",
			"store food product",
		}
		for _, q := range query {
			s, err := testUserRepo.Find(q, nil, nil, 0, 0)
			if err != nil || len(s) < 1 || s[0].Name != testUser.Name {
				t.Error(err)
			}
		}
	})
}

func TestUserRepo_UpdateOne(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testUserRepo.UpdateOne(user.User{
			ID:   testUserID,
			Name: "Test User Change name",
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestUserRepo_DeleteOne(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testUserRepo.DeleteOne(testUserID)
		if err != nil {
			t.Error(err)
		}
	})
}
