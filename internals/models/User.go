package models

import (
	"time"
	"wechat-back/internals/generators"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CONSTANTS
const (
	CODE_NOT_VALID = 46
	CODE_VALID     = 65
)

// User structure that onlines the user
type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	Email         string             `json:"email" bson:"email"`
	ValidCode     int                `json:"valid_code" bson:"valid_code"`
	TempCode      int                `json:"temp_code" bson:"temp_code"`
	CodeTimestamp time.Time          `json:"code_timestamp" bson:"code_timestamp"`
	Credentials   UserCredentials    `json:"credentials" bson:"crendentials"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}

/*
UserCredentials holds necesary
information about the user to get
full app functionality
*/
type UserCredentials struct {
	PushToken string `json:"push_token" bson:"push_token"`
}

/*
FormatUserModel
fill the fields the user is
not allowed to fill and returns
the formated model
*/
func FormatUserModel(m *User) *User {

	m.ID = primitive.NewObjectID()
	m.ValidCode = CODE_VALID
	m.TempCode, _ = generators.Generate6DigitCode()
	m.CodeTimestamp = time.Now()
	m.CreatedAt = time.Now()

	return m
}
