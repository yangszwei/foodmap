package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User user model
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty" validate:"required,max=50"`
	Email     string             `json:"email" bson:"name,omitempty" validate:"required,email"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"upd,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"cre,omitempty"`
}
