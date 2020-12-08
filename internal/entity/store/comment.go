package store

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"time"
)

// Comment user comment
type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"uid" validate:"required"`
	Stars     uint8              `json:"stars" bson:"star" validate:"required,min=0,max=5"`
	Message   string             `json:"message,omitempty" bson:"msg" validate:"max=200"`
	IPAddr    net.IP             `json:"ip_addr,omitempty" bson:"ip" validate:"required"`
	UserAgent string             `json:"user_agent,omitempty" bson:"ua" validate:"required"`
	CreatedAt *time.Time         `json:"created_at" bson:"cre" validate:"required"`
}
