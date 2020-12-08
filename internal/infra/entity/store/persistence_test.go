package store_test

import (
	"foodmap/internal/entity/store"
	"foodmap/internal/infra/config"
	infra "foodmap/internal/infra/entity/store"
	"foodmap/internal/infra/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net"
	"testing"
	"time"
)

var (
	testStore = store.Store{
		Name:          "Test Store",
		Description:   "$set:{\"test\": 1}",
		BusinessHours: []store.BusinessHoursRule{
			{
				FromDay:  time.Monday,
				ToDay:    time.Friday,
				FromTime: "08:00",
				ToTime:   "22:00",
			},
		},
		Categories:    []string{ "food" },
		PriceLevel:    store.PriceCheap,
		Menu:          []store.Product{
			{
				Name:        "Test Product",
				Description: "Test description",
				Category:    "Main",
				Variants:    []store.Variant{
					{
						Name:  "大",
						Price: 50,
					},
					{
						Name: "小",
						Price: 40,
					},
				},
			},
		},
	}
	testStoreID primitive.ObjectID
	testStoreCommentID primitive.ObjectID
	testStoreRepo *infra.Repo
)

func TestNewStoreRepo(t *testing.T) {
	cfg, err := config.Open("../../../../.env")
	if err != nil {
		log.Fatalln(err)
	}
	db, err := persistence.Open(cfg.DB)
	if err != nil {
		log.Fatalln(err)
	}
	testStoreRepo, err = infra.NewRepo(&db)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestStoreRepo_InsertOne(t *testing.T) {
	var err error
	testStoreID, err = testStoreRepo.InsertOne(testStore)
	if err != nil {
		t.Error(err)
	}
}

func TestStoreRepo_FindOneByID(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		s, err := testStoreRepo.FindOneByID(testStoreID, bson.M{ "name": 1 })
		if err != nil {
			t.Error(err)
		}
		if s.Name != testStore.Name {
			t.Error(s)
		}
		t.Log(s)
	})
}

func TestStoreRepo_Find(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		query := []string{
			"store",
			"food",
			"Product",
			"Main",
			"store food product",
		}
		for _, q := range query {
			s, err := testStoreRepo.Find(q, nil, nil, 0, 0)
			if err != nil || len(s) < 1 || s[0].Name != testStore.Name {
				t.Error(err)
			}
		}
	})
}

func TestStoreRepo_UpdateOne(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testStoreRepo.UpdateOne(store.Store{
			ID: testStoreID,
			Description: "updated description",
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestStoreRepo_InsertComment(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		_, err := testStoreRepo.InsertComment(testStoreID, store.Comment{
			UserID:    primitive.NewObjectID(),
			Stars:     4,
			Message:   "test comment",
			IPAddr:    net.ParseIP("192.0.2.1"),
			UserAgent: "test user agent",
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestStoreRepo_FindComments(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		c, err := testStoreRepo.FindComments(testStoreID, 10, 0)
		if err != nil || len(c) < 1 {
			t.Fatal(err)
		}
		testStoreCommentID = c[0].ID
	})
}

func TestStoreRepo_DeleteComment(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testStoreRepo.DeleteComment(testStoreID, testStoreCommentID)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestStoreRepo_DeleteOne(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testStoreRepo.DeleteOne(testStoreID)
		if err != nil {
			t.Error(err)
		}
	})
}
