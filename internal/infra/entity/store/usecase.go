package store

import (
	"encoding/json"
	"foodmap/internal/entity/store"
	"foodmap/internal/infra/errors"
	"foodmap/internal/infra/object"
	"foodmap/internal/infra/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"net"
	"reflect"
	"strings"
	"time"
)

// NewUsecase setup and return an instance of Usecase
func NewUsecase(r store.IStoreRepo, v *validator.Validator) *Usecase {
	s := new(Usecase)
	s.r = r
	s.v = v
	return s
}

// Usecase implement store.IStoreUsecase
type Usecase struct {
	r store.IStoreRepo
	v *validator.Validator
}

// Create validate and insert a store record to repository
func (u Usecase) CreateOne(document object.H) (string, error) {
	doc, err := processStoreDocumentInput(document)
	if err != nil {
		return "", err
	}
	record, err := toStoreEntity(doc)
	if err != nil {
		return "", err
	}
	if err := u.v.Validate(record) ; err != nil {
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
	if list := isFieldsValid(fields) ; list != nil {
		return nil, errors.New("invalid fields", strings.Join(list, ","))
	}
	record, err := u.r.FindOneByID(storeID, toProjection(fields))
	if err != nil {
		return nil, err
	}
	return processStoreDocumentOutput(record), nil
}

// Find find a list of store records from repository, fields should be a
// comma-separated list of available fields (refer to store.Store's json tag)
func (u Usecase) Find(query, categories, fields string, limit, skip int64) ([]object.H, error) {
	if f := isFieldsValid(fields) ; f != nil {
		return []object.H{{ "fields": f }}, errors.New("invalid fields")
	}
	records, err := u.r.Find(query, split(categories), toProjection(fields), limit, skip)
	if err != nil {
		return nil, err
	}
	result := []object.H{}
	for _, record := range records {
		result = append(result, processStoreDocumentOutput(record))
	}
	return result, nil
}

// UpdateOne update a store record on repository
func (u Usecase) UpdateOne(document object.H) error {
	var update, compare store.Store
	id, ok := document["id"].(string)
	if !ok {
		return errors.NewValidationError("id", "required")
	}
	storeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewValidationError("id", "invalid")
	}
	if update, err = toStoreEntity(document) ; err != nil {
		return err
	}
	if compare, err = u.r.FindOneByID(storeID, nil) ; err != nil {
		return err
	}
	m, _ := json.Marshal(document)
	_ = json.Unmarshal(m, &compare)
	if err := u.v.Validate(compare) ; err != nil {
		return err
	}
	return u.r.UpdateOne(update)
}

// DeleteOne delete a store record from repository
func (u Usecase) DeleteOne(id string) error {
	storeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return u.r.DeleteOne(storeID)
}

// CreateComment add a comment to store record
func (u Usecase) CreateComment(storeID string, comment object.H) (string, error) {
	id, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return "", err
	}
	if d, ok := comment["ip_addr"].(string) ; ok {
		comment["ip_addr"] = net.ParseIP(d)
	}
	record, err := toStoreCommentEntity(comment)
	if err != nil {
		return "", err
	}
	t := time.Now()
	record.CreatedAt = &t
	if err := u.v.Validate(record) ; err != nil {
		return "", err
	}
	if id, err = u.r.InsertComment(id, record) ; err != nil {
		return "", err
	}
	return id.Hex(), err
}

// FindComments find a list of comments in store record
func (u Usecase) FindComments(storeID string, admin bool, limit, skip int64) ([]object.H, error) {
	id, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return nil, err
	}
	records, err := u.r.FindComments(id, limit, skip)
	if err != nil {
		return nil, err
	}
	result := []object.H{}
	for _, record := range records {
		doc := toStoreCommentDocument(record)
		doc["id"] = record.ID.Hex()
		doc["user_id"] = record.UserID.Hex()
		if _, exist := doc["ip_addr"] ; exist && !admin {
			delete(doc, "ip_addr")
		}
		if _, exist := doc["user_agent"] ; exist && !admin {
			delete(doc, "user_agent")
		}
		result = append(result, doc)
	}
	return result, nil
}

// DeleteComment remove a comment from store record
func (u Usecase) DeleteComment(storeID, commentID string) error {
	sid, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return errors.NewValidationError("store_id", "invalid")
	}
	cid, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.NewValidationError("comment_id", "invalid")
	}
	return u.r.DeleteComment(sid, cid)
}

// toStoreEntity convert object.H to store.Store
func toStoreEntity(doc object.H) (store.Store, error) {
	var result store.Store
	data, err := json.Marshal(doc)
	if err != nil {
		return store.Store{}, err
	}
	_ = json.Unmarshal(data, &result)
	return result, nil
}

// toStoreDocument convert store.Store to object.H
func toStoreDocument(record store.Store) object.H {
	var result object.H
	data, _ := json.Marshal(record)
	_ = json.Unmarshal(data, &result)
	return result
}

// toStoreCommentEntity convert object.H to store.Comment
func toStoreCommentEntity(doc object.H) (store.Comment, error) {
	var result store.Comment
	data, err := json.Marshal(doc)
	if err != nil {
		return store.Comment{}, err
	}
	_ = json.Unmarshal(data, &result)
	return result, nil
}

// toStoreCommentDocument convert store.Comment to object.H
func toStoreCommentDocument(record store.Comment) object.H {
	var result object.H
	data, _ := json.Marshal(record)
	_ = json.Unmarshal(data, &result)
	return result
}

func processStoreDocumentInput(doc object.H) (object.H, error) {
	if d, exist := doc["business_hours"].([]interface{}) ; exist {
		bh, err := processStoreBusinessHoursInput(d)
		if err != nil {
			return nil, err
		}
		doc["business_hours"] = bh
	}
	if d, exist := doc["price_level"] ; exist {
		pl, err := processStorePriceLevelInput(d)
		if err != nil {
			return nil, err
		}
		doc["price_level"] = pl
	}
	return doc, nil
}

func processStoreBusinessHoursInput(bhs []interface{}) ([]object.H, error) {
	var result []object.H
	for _, bh := range bhs {
		doc := bh.(map[string]interface{})
		var ( day, time []interface{} ; ok bool )
		if day, ok = doc["day"].([]interface{}) ; !ok || len(day) != 2 {
			return nil, errors.NewValidationError("business_hours.day", "invalid")
		}
		if time, ok = doc["time"].([]interface{}) ; !ok || len(time) != 2 {
			return nil, errors.NewValidationError("business_hours.time", "invalid")
		}
		obj := object.H{
			"from_day": day[0],
			"to_day": day[1],
			"from_time": time[0],
			"to_time": time[1],
		}
		data, _ := json.Marshal(obj)
		var record store.BusinessHoursRule
		if err := json.Unmarshal(data, &record) ; err != nil ||
			record.FromDay == 0 || record.ToDay == 0 || record.FromTime == "" || record.ToTime == "" {
			return nil, errors.NewValidationError("business_hours", "invalid")
		}
		result = append(result, obj)
	}
	return result, nil
}

func processStorePriceLevelInput(pl interface{}) (rune, error) {
	pl, ok := pl.(string)
	if !ok {
		return 0, nil
	}
	switch pl {
	case "cheap":
		return store.PriceCheap, nil
	case "medium":
		return store.PriceMedium, nil
	case "expensive":
		return store.PriceExpensive, nil
	default:
		return store.PriceUnknown, errors.NewValidationError("price_level", "invalid")
	}
}

func processStoreDocumentOutput(record store.Store) object.H {
	doc := toStoreDocument(record)
	if _, exist := doc["id"] ; exist {
		doc["id"] = record.ID
	}
	if _, exist := doc["business_hours"] ; exist {
		doc["business_hours"] = processStoreBusinessHoursOutput(record.BusinessHours)
	}
	return toStoreDocument(record)
}

func processStoreBusinessHoursOutput(bhs []store.BusinessHoursRule) [7][][2]string {
	var result [7][][2]string
	for _, bh := range bhs {
		min := math.Min(float64(bh.FromDay), float64(bh.ToDay))
		max := math.Max(float64(bh.FromDay), float64(bh.ToDay))
		for i := int(min) ; i <= int(max) ; i++ {
			result[i - 1] = append(result[i - 1], [2]string{ bh.FromTime, bh.ToTime })
		}
	}
	return result
}

// isFieldsValid returns a list of unrecognized fields
func isFieldsValid(fields string) (result []string) {
	tags := split(fields)
	if len(tags) == 0 {
		return nil
	}
	validTags := []string{
		"id", "name", "description", "is_open", "business_hours", "categories",
		"price_level", "menu", "average_stars", "updated_at", "created_at",
	}
	for _, tag := range tags {
		var exist bool
		if !exist {
			for _, t := range validTags {
				if tag == t {
					exist = true
					break
				}
			}
		}
		if !exist {
			result = append(result, tag)
		}
	}
	return
}

// toProjection convert fields query to projection document
func toProjection(fields string) bson.M {
	var result = make(bson.M)
	tags := split(fields)
	if len(tags) == 0 {
		// do not request comments by default
		return bson.M{ "cmnt": false }
	}
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

// split split string by comma and remove empty ones
func split(str string) (result []string) {
	for _, i := range strings.Split(str, ",") {
		if i != "" {
			result = append(result, i)
		}
	}
	return
}
