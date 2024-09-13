package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
MESSAGE TYPES
*/
const (

	// MESSAGE_STANDARD_TYPE
	// Type of body which the payload is a text
	MESSAGE_STANDARD_TYPE = 25

	// MESSAGE_MEDIA_TYPE
	// Type of body which the payload is a video or an image
	MESSAGE_MEDIA_TYPE = 26

	// MESSAGE_FILE_TYPE
	// type of body which the payload is an attachment
	MESSAGE_FILE_TYPE = 32
)

/*
GroupChatLog
// Represent a basic message structure
*/
type GroupChatLog struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	TargetID   primitive.ObjectID `json:"target" bson:"target"`
	AuthorID   primitive.ObjectID `json:"author_id" bson:"author_id"`
	ContentID  string             `json:"content_id" bson:"content_id"`
	AuthorName string             `json:"author_name" bson:"author_name"`
	BodyType   int                `json:"body_type" bson:"body_type"`
	Body       string             `json:"body" bson:"body"`
	Alt        string             `json:"alt" bson:"alt"`
	Created_at time.Time          `json:"created_at" bson:"created_at"`
}
