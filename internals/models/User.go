package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User structure that onlines the user
type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	ValidCode int                `json:"valid_code" bson:"valid_code"`
	TempCode  int                `json:"temp_code"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
