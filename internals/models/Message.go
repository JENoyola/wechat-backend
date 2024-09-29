package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
CONSTANTS
*/
const (
	MESSAGE_EDTIED = 12

	// MESSAGE_UNEDITED
	// Means that the message has not been edited
	MESSAGE_UNEDITED = 45
)

// GroupChatTextLog Represent a message structure for groups
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

// GroupChatContentLog content message structure for groups
type GroupChatContentLog struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	TargetID     primitive.ObjectID `json:"target_id" bson:"target_id"`
	AuthorID     primitive.ObjectID `json:"author_id" bson:"author_id"`
	ContentID    string             `json:"content_id" bson:"content_id"`
	AuthorName   string             `json:"author_name" bson:"author_name"`
	BodyType     int                `json:"body_type" bson:"body_type"`
	Body         string             `json:"body" bson:"body"`
	Media        []string           `json:"media" bson:"media"`
	Placeholders []string           `json:"placeholders" bson:"placeholders"`
	Created_at   time.Time          `json:"created_at" bson:"created_at"`
}

func (g *GroupChatTextLog) FormatTextChatLog(groupID, author_id primitive.ObjectID, authorname, body string) {
	g.ID = primitive.NewObjectID()
	g.TargetID = groupID
	g.AuthorID = author_id
	g.ContentID = "N/A"
	g.AuthorName = authorname
	g.BodyType = MESSAGE_TYPE_TEXT
	g.Body = body
	g.Alt = ""
	g.Created_At = time.Now()
}

func (p *GroupChatContentLog) FormatContentChatLog(targetID, author primitive.ObjectID, authorName, body, contentID string, files []string, placeholders []string, MessageType int) {
	p.ID = primitive.NewObjectID()
	p.TargetID = targetID
	p.AuthorID = author
	p.ContentID = contentID
	p.AuthorName = authorName
	p.BodyType = MessageType
	p.Body = body
	p.Media = files
	p.Placeholders = placeholders
	p.Created_at = time.Now()
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

type P2PContentChatLog struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	TargetID     primitive.ObjectID `json:"target_id" bson:"target_id"`
	AuthorID     primitive.ObjectID `json:"author_id" bson:"author_id"`
	ContentID    string             `json:"content_id" bson:"content_id"`
	AuthorName   string             `json:"author_name" bson:"author_name"`
	BodyType     int                `json:"body_type" bson:"body_type"`
	Body         string             `json:"body" bson:"body"`
	Media        []string           `json:"media" bson:"media"`
	Placeholders []string           `json:"placeholders" bson:"placeholders"`
	Created_at   time.Time          `json:"created_at" bson:"created_at"`
}

// FormatContentChatLog fills fields on chatlogs that contains media
func (p *P2PContentChatLog) FormatContentChatLog(targetID, author primitive.ObjectID, authorName, body, contentID string, files []string, placeholders []string, MessageType int) {
	p.ID = primitive.NewObjectID()
	p.TargetID = targetID
	p.AuthorID = author
	p.ContentID = contentID
	p.AuthorName = authorName
	p.BodyType = MessageType
	p.Body = body
	p.Media = files
	p.Placeholders = placeholders
	p.Created_at = time.Now()
}

// FormatTextLog formats P2PTextChatLog
func (p *P2PTextChatLog) FormatTextLog(targetID, author primitive.ObjectID, authorName, body string) {
	p.ID = primitive.NewObjectID()
	p.TargetID = targetID
	p.AuthorID = author
	p.ContentID = "N/A"
	p.AuthorName = authorName
	p.BodyType = MESSAGE_TYPE_TEXT
	p.Body = body
	p.Created_at = time.Now()
}
