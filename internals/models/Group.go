package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id"`
	Name         string               `json:"name" bson:"name"`
	Description  string               `json:"description" bson:"description"`
	Participants []primitive.ObjectID `json:"participants" bson:"participants"`
	Admins       []primitive.ObjectID `json:"admins" bson:"admins"`
	CreatedAt    time.Time            `json:"created_at" bson:"created_at"`
}
