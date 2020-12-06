package usecase_test

import (
	"encoding/json"
	"foodmap/internal/infra/object"
	"foodmap/internal/infra/persistence"
	"foodmap/internal/infra/validator"
	"foodmap/internal/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"testing"
)

var (
	testStoreID, testStoreCommentID string
	testStoreUsecase *usecase.StoreUsecase
)

func TestNewStoreUsecase(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		db, err := Open()
		if err != nil {
			log.Fatal(err)
		}
		r, err := persistence.NewStoreRepo(&db)
		if err != nil {
			log.Fatal(err)
		}
		testStoreUsecase = usecase.NewStoreUsecase(r, validator.New())
	})
}

func TestStoreUsecase_CreateOne(t *testing.T) {
	t.Run("should success", func(t *testing.T) {
		id, err := testStoreUsecase.CreateOne(object.H{
			"name": "Test Store",
			"description": "Test description",
			"business_hours": []object.H{
				{
					"day": []int{ 1, 5 },
					"time": []string{ "08:00", "22:00" },
				},
			},
			"categories": []string{ "a", "b" },
			"price_level": "medium",
			"menu": []object.H{
				{
					"name": "Test Product",
					"description": "test description",
					"category": "main",
					"price": 200,
				},
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
		testStoreID = id
	})
}

func TestStoreUsecase_FindOneByID(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		data, err := testStoreUsecase.FindOneByID(testStoreID, "name,categories")
		if err != nil {
			t.Fatal(err)
		}
		var n []string
		d := data["categories"].([]interface{})
		dat, _ := json.Marshal(d)
		_ = json.Unmarshal(dat, &n)
		if data["name"].(string) != "Test Store" || strings.Join(n, ",") != "a,b" {
			t.Error(data)
		}
	})
}

func TestStoreUsecase_Find(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		result, err := testStoreUsecase.Find("test store", nil, "name", 10, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(result) < 1 || result[0]["name"].(string) != "Test Store" {
			t.Error(result)
		}
	})
}

func TestStoreUsecase_UpdateOne(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testStoreUsecase.UpdateOne(object.H{
			"id": testStoreID,
			"description": "updated description",
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestStoreUsecase_CreateComment(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testStoreUsecase.CreateComment(testStoreID, object.H{
			"user_id": primitive.NewObjectID().Hex(),
			"stars": 4,
			"ip_addr": "127.0.0.1",
			"user_agent": "test ua",
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestStoreUsecase_FindComments(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		records, err := testStoreUsecase.FindComments(testStoreID, 10, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(records) < 1 {
			t.Fatal(records)
		}
		type cmnt struct {
			Stars int
		}
		var c cmnt
		d, _ := json.Marshal(records[0])
		_ = json.Unmarshal(d, &c)
		if c.Stars != 4 {
			t.Fatal(records)
		}
		testStoreCommentID = records[0]["id"].(string)
	})
}

func TestStoreUsecase_DeleteComment(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testStoreUsecase.DeleteComment(testStoreID, testStoreCommentID)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestStoreUsecase_DeleteOne(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		err := testStoreUsecase.DeleteOne(testStoreID)
		if err != nil {
			t.Error(err)
		}
	})
}
