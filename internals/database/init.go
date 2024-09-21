package database

import (
	"context"
	"errors"
	"log"
	"os"
	"wechat-back/internals/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBHUB interface {
	// users
	FindUserDB(string) (models.User, bool, error)
	InsertUserDB(models.User) (string, error)
	UpdateUserAccountDB(map[string]any, string) error
	GetUsers(int, string) ([]*models.User, error)

	// groups
	GetGroupDB(string) (*models.Group, error)
	InsertGroupDB(models.Group) (string, error)
	UpdateGroupDB(map[string]any, primitive.ObjectID) error
	DeleteGroupDB(string) error
	SearchGroups(int, string) ([]*models.Group, error)
}

// ERRORS
var (
	ErrNoModified = errors.New("no documents modified")
	ErrNoDeleted  = errors.New("no deleted")
)

type DB struct {
	Client   *mongo.Client
	Database string
}

// StartDatabase makes the connection with database and returns the connection
func StartDatabase() *DB {
	ctx := context.TODO()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_URI")))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &DB{
		Client:   c,
		Database: os.Getenv("DB_DATABASE"),
	}
}

// FormatUserCollection Formats the collection for users
func (db *DB) FormatUserCollection() *mongo.Collection {
	return db.Client.Database(db.Database).Collection(os.Getenv("DB_USERS"))
}

// FormatGroupCollection Formats the collection for groupd
func (db *DB) FormatGroupCollection() *mongo.Collection {
	return db.Client.Database(db.Database).Collection(os.Getenv("DB_GROUPD"))
}

// FormatUserChatlogs Formats the collection for user chatlogs
func (db *DB) FormatUserChatlogs() *mongo.Collection {
	return db.Client.Database(db.Database).Collection(os.Getenv("DB_USR_CHLOGS"))
}

// FormatGroupChatlogs Formats the collection for group chat logs
func (db *DB) FormatGroupChatlogs() *mongo.Collection {
	return db.Client.Database(db.Database).Collection(os.Getenv("DB_GR_CHLOGS"))
}
