package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
CONSTANTS
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

	// MESSAGE_EDTIED
	// Means that the message has been edited by the author
	MESSAGE_EDTIED = 12

	// MESSAGE_UNEDITED
	// Means that the message has not been edited
	MESSAGE_UNEDITED = 45
)

/*
GroupChatLog
// Represent a message structure for groups
*/
type GroupChatTextLog struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	TargetID   primitive.ObjectID `json:"target_id" bson:"target_id"`
	AuthorID   primitive.ObjectID `json:"author_id" bson:"author_id"`
	ContentID  string             `json:"content_id" bson:"content_id"`
	AuthorName string             `json:"author_name" bson:"author_name"`
	BodyType   int                `json:"body_type" bson:"body_type"`
	Body       string             `json:"body" bson:"body"`
	Alt        string             `json:"alt" bson:"alt"`
	Created_At time.Time          `json:"created_at" bson:"created_at"`
}

func (g *GroupChatTextLog) FormatTextChatLog(groupID, author_id primitive.ObjectID, authorname, body string) {
	g.ID = primitive.NewObjectID()
	g.TargetID = groupID
	g.AuthorID = author_id
	g.ContentID = "N/A"
	g.AuthorName = authorname
	g.BodyType = MESSAGE_STANDARD_TYPE
	g.Body = body
	g.Alt = ""
	g.Created_At = time.Now()
}

/*
P2PChatLog
Represents a message structure for private conversations
*/
type P2PTextChatLog struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	TargetID   primitive.ObjectID `json:"target_id" bson:"target_id"`
	AuthorID   primitive.ObjectID `json:"author_id" bson:"author_id"`
	ContentID  string             `json:"content_id" bson:"content_id"`
	AuthorName string             `json:"author_name" bson:"author_name"`
	BodyType   int                `json:"body_type" bson:"body_type"`
	Body       string             `json:"body" bson:"body"`
	Created_at time.Time          `json:"created_at" bson:"created_at"`
}

func (p *P2PTextChatLog) FormatTextLog(targetID, author primitive.ObjectID, authorName, body string) {
	p.ID = primitive.NewObjectID()
	p.TargetID = targetID
	p.AuthorID = author
	p.ContentID = "N/A"
	p.AuthorName = authorName
	p.BodyType = MESSAGE_STANDARD_TYPE
	p.Body = body
	p.Created_at = time.Now()
}
