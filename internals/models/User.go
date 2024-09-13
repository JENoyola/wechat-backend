package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CONSTANTS
const (
	CODE_NOT_VALID = 46
	CODE_VALID     = 65
)

// User structure that onlines the user
type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	ValidCode int                `json:"valid_code" bson:"valid_code"`
	TempCode  int                `json:"temp_code" bson:"temp_code"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
