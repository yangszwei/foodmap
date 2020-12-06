package usecase

import (
	"encoding/json"
	"foodmap/internal/entity/store"
	"foodmap/internal/infra/errors"
	"foodmap/internal/infra/object"
	"foodmap/internal/infra/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"reflect"
	"strings"
)

// NewStoreUsecase setup and return an instance of StoreUsecase
func NewStoreUsecase(r store.IStoreRepo, v *validator.Validator) *StoreUsecase {
	s := new(StoreUsecase)
	s.r = r
	s.v = v
	return s
}

// StoreUsecase implement store.IStoreUsecase
type StoreUsecase struct {
	r store.IStoreRepo
	v *validator.Validator
}

// Create validate and insert a store record to repository
func (s StoreUsecase) CreateOne(document object.H) (string, error) {
	var record store.Store
	if doc, exist := document["business_hours"].([]object.H) ; exist {
		bh, err := toRepoBusinessHours(doc)
		if err != nil {
			return "", err
		}
		document["business_hours"] = bh
	}
	pl, err := toRepoPriceLevel(document["price_level"].(string))
	if err != nil {
		return "", err
	}
	document["price_level"] = pl
	data, _ := json.Marshal(document)
	_ = json.Unmarshal(data, &record)
	if err := s.v.Validate(record) ; err != nil {
		return "", err
	}
	id, err := s.r.InsertOne(record)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

// FindOneByID find a store record by ID from repository with selected fields
// fields should be a comma-separated list of available fields (refer to
// store.Store's json tag)
func (s StoreUsecase) FindOneByID(id string, fields string) (object.H, error) {
	var result object.H
	storeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	record, err := s.r.FindOneByID(storeID, toProjection(fields))
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(record)
	_ = json.Unmarshal(data, &result)
	return result, nil
}

// Find find a list of store records from repository, fields should be a
// comma-separated list of available fields (refer to store.Store's json tag)
func (s StoreUsecase) Find(query string, filter object.H, fields string, limit, skip int) ([]object.H, error) {
	var ( f store.Store ; result []object.H )
	data, _ := json.Marshal(filter)
	_ = json.Unmarshal(data, &f)
	records, err := s.r.Find(query, f, toProjection(fields), int64(limit), int64(skip))
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		var res object.H
		data, _ := json.Marshal(record)
		_ = json.Unmarshal(data, &res)
		result = append(result, res)
	}
	return result, nil
}

// UpdateOne update a store record on repository
func (s StoreUsecase) UpdateOne(document object.H) error {
	var local, record store.Store
	id, ok := document["id"].(string)
	if !ok {
		return errors.NewValidationError("id", "required")
	}
	storeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	local, err = s.r.FindOneByID(storeID, nil)
	if err != nil {
		return err
	}
	if doc, exist := document["business_hours"].([]object.H) ; exist {
		bh, err := toRepoBusinessHours(doc)
		if err != nil {
			return err
		}
		document["business_hours"] = bh
	}
	if doc, exist := document["price_level"].(string) ; exist {
		pl, err := toRepoPriceLevel(doc)
		if err != nil {
			return err
		}
		document["price_level"] = pl
	}
	data, _ := json.Marshal(document)
	_ = json.Unmarshal(data, &local)
	if err := s.v.Validate(local) ; err != nil {
		return err
	}
	_ = json.Unmarshal(data, &record)
	return s.r.UpdateOne(record)
}

// DeleteOne delete a store record from repository
func (s StoreUsecase) DeleteOne(id string) error {
	storeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.r.DeleteOne(storeID)
}

// CreateComment add a comment to store record
func (s StoreUsecase) CreateComment(storeID string, comment object.H) error {
	id, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return err
	}
	comment["ip_addr"] = net.ParseIP(comment["ip_addr"].(string))
	var record store.Comment
	data, _ := json.Marshal(comment)
	_ = json.Unmarshal(data, &record)
	return s.r.InsertComment(id, record)
}

// FindComments find a list of comments in store record
func (s StoreUsecase) FindComments(storeID string, limit, skip int) ([]object.H, error) {
	id, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return nil, err
	}
	records, err := s.r.FindComments(id, int64(limit), int64(skip))
	if err != nil {
		return nil, err
	}
	var result []object.H
	for _, record := range records {
		var res object.H
		data, _ := json.Marshal(record)
		_ = json.Unmarshal(data, &res)
		result = append(result, res)
	}
	return result, nil
}

// DeleteComment remove a comment from store record
func (s StoreUsecase) DeleteComment(storeID, commentID string) error {
	sid, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return err
	}
	cid, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return err
	}
	return s.r.DeleteComment(sid, cid)
}

// toProjection convert a comma-separated list of field names to projection
// document
func toProjection(fields string) bson.M {
	var result = make(bson.M)
	tags := strings.Split(fields, ",")
	rt := reflect.TypeOf(store.Store{})
	for i := 0 ; i < rt.NumField() ; i++ {
		f := rt.Field(i)
		for _, tag := range tags {
			if tag == strings.Split(f.Tag.Get("json"), ",")[0] {
				result[strings.Split(f.Tag.Get("bson"), ",")[0]] = true
			}
		}
	}
	return result
}

// convert input business hours document to repository format
func toRepoBusinessHours(bhs []object.H) ([]object.H, error) {
	var result []object.H
	for _, bh := range bhs {
		day := bh["day"].([]int)
		if len(day) != 2 {
			return nil, errors.NewValidationError("business_hours.day", "invalid")
		}
		time := bh["time"].([]string)
		if len(time) != 2 {
			return nil, errors.NewValidationError("business_hours.day", "invalid")
		}
		result = append(result, object.H{
			"from_day": day[0],
			"to_day": day[1],
			"from_time": time[0],
			"to_time": time[1],
		})
	}
	return result, nil
}

// toRepoPriceLevel convert price level to repo format
func toRepoPriceLevel(pl string) (rune, error) {
	switch pl {
	case "cheap":
		return store.PriceCheap, nil
	case "medium":
		return store.PriceMedium, nil
	case "expensive":
		return store.PriceExpensive, nil
	default:
		return 0, errors.NewValidationError("price_level", "invalid")
	}
}
