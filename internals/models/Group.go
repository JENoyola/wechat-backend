package models

import (
	"time"
	"wechat-back/internals/generators"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group basic group structure
type Group struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id"`
	GroupID      string               `json:"group_id" bson:"group_id"`
	Name         string               `json:"name" bson:"name"`
	Description  string               `json:"description" bson:"description"`
	Participants []primitive.ObjectID `json:"participants" bson:"participants"`
	Admins       []primitive.ObjectID `json:"admins" bson:"admins"`
	ProfileImage string               `json:"profile_image" bson:"profile_image"`
	CreatedAt    time.Time            `json:"created_at" bson:"created_at"`
}

// FormatGroup adds the necessary information to the structure
func FormatGroup(g *Group) *Group {
	g.ID = primitive.NewObjectID()
	g.GroupID = generators.GenerateUniqueID(g.ID.Hex(), g.Name)
	g.CreatedAt = time.Now()
	return g
}
