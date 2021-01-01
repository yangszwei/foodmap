package user

import (
	"encoding/json"
	"foodmap/internal/entity/user"
	"foodmap/internal/infra/delivery"
	"foodmap/internal/infra/errors"
	"foodmap/internal/infra/object"
	"foodmap/internal/infra/validator"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// validTags provide a string contained all valid tags
func validTags() (tags []string) {
	tags = []string{
		"id", "name", "description", "is_open", "business_hours", "categories",
		"price_level", "menu", "average_stars", "updated_at", "created_at",
	}
	return
}

// NewUsecase setup and return an instance of Usecase
func NewUsecase(r user.IUserRepo, v *validator.Validator) *Usecase {
	s := new(Usecase)
	s.r = r
	s.v = v
	return s
}

// Usecase implement user.IUserRepo
type Usecase struct {
	r user.IUserRepo
	v *validator.Validator
}

// Create validate and insert a user record to repository
func (u Usecase) CreateOne(document object.H) (string, error) {
	doc, err := processUserDocumentInput(document)
	if err != nil {
		return "", err
	}
	record, err := toUserEntity(doc)
	if err != nil {
		return "", err
	}
	if err := u.v.Validate(record); err != nil {
		return "", err
	}
	id, err := u.r.InsertOne(record)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

// FindOneByID find a store record by ID from repository with selected fields
// fields should be a comma-separated list of available fields (refer to
// store.Store's json tag)
func (u Usecase) FindOneByID(id string, fields string) (object.H, error) {
	storeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if list := validator.IsFieldsValid(validTags(), fields); list != nil {
		return nil, errors.New("invalid fields", strings.Join(list, ","))
	}
	record, err := u.r.FindOneByID(storeID, toProjection(fields))
	if err != nil {
		return nil, err
	}
	return processUserDocumentOutput(record), nil
}

// Find find a list of store records from repository, fields should be a
// comma-separated list of available fields (refer to store.Store's json tag)
func (u Usecase) Find(query, fields string, limit, skip int64) ([]object.H, error) {
	if f := validator.IsFieldsValid(validTags(), fields); f != nil {
		return []object.H{{"fields": f}}, errors.New("invalid fields")
	}
	records, err := u.r.Find(query, toProjection(fields), limit, skip)
	if err != nil {
		return nil, err
	}
	result := []object.H{}
	for _, record := range records {
		result = append(result, processUserDocumentOutput(record))
	}
	return result, nil
}

// UpdateOne update a store record on repository
func (u Usecase) UpdateOne(document object.H) error {
	var update, compare user.User
	id, ok := document["id"].(string)
	if !ok {
		return errors.NewValidationError("id", "required")
	}
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewValidationError("id", "invalid")
	}
	if update, err = toUserEntity(document); err != nil {
		return err
	}
	if compare, err = u.r.FindOneByID(userID, nil); err != nil {
		return err
	}
	m, _ := json.Marshal(document)
	_ = json.Unmarshal(m, &compare)
	if err := u.v.Validate(compare); err != nil {
		return err
	}
	return u.r.UpdateOne(update)
}

// DeleteOne delete a store record from repository
func (u Usecase) DeleteOne(id string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return u.r.DeleteOne(userID)
}

// toStoreEntity convert object.H to store.Store
func toUserEntity(doc object.H) (user.User, error) {
	var result user.User
	data, err := json.Marshal(doc)
	if err != nil {
		return user.User{}, err
	}
	_ = json.Unmarshal(data, &result)
	return result, nil
}

// toStoreDocument convert store.Store to object.H
func toUserDocument(record user.User) object.H {
	var result object.H
	data, _ := json.Marshal(record)
	_ = json.Unmarshal(data, &result)
	return result
}

func processUserDocumentInput(doc object.H) (object.H, error) {
	return doc, nil
}

func processUserDocumentOutput(record user.User) object.H {
	doc := toUserDocument(record)
	return doc
}

// toProjection convert fields query to projection document
func toProjection(fields string) bson.M {
	var result = make(bson.M)
	tags := delivery.Split(fields)
	if len(tags) == 0 {
		// do not request comments by default
		return bson.M{"cmnt": false}
	}
	rt := reflect.TypeOf(user.User{})
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		for _, tag := range tags {
			if tag == strings.Split(f.Tag.Get("json"), ",")[0] {
				result[strings.Split(f.Tag.Get("bson"), ",")[0]] = true
			}
		}
	}
	return result
}
