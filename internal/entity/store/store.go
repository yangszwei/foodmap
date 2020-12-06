package store

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Store store model
type Store struct {
	ID            primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string              `json:"name,omitempty" bson:"name,omitempty" validate:"required,max=50"`
	Description   string              `json:"description,omitempty" bson:"desc,omitempty" validate:"max=1000"`
	BusinessHours []BusinessHoursRule `json:"business_hours,omitempty" bson:"bh,omitempty" validate:"required"`
	Categories    []string            `json:"categories,omitempty" bson:"cat,omitempty"`
	PriceLevel    rune                `json:"price_level,omitempty" bson:"pl,omitempty" validate:"required"`
	Menu          []Product           `json:"menu,omitempty" bson:"menu,omitempty"`
	Comments      []Comment           `json:"comments,omitempty" bson:"cmnt,omitempty"`
	UpdatedAt     *time.Time          `json:"updated_at,omitempty" bson:"upd,omitempty"`
	CreatedAt     *time.Time          `json:"created_at,omitempty" bson:"cre,omitempty"`
}

// Price levels
const (
	PriceCheap = 'c'
	PriceMedium = 'm'
	PriceExpensive = 'e'
)
