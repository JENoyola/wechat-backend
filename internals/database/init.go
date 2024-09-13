package database

import (
	"context"
	"errors"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ERRORS
var (
	ErrNoModified = errors.New("no documents modified")
	ErrNoDeleted  = errors.New("no deleted")
)

type DB struct {
	client   *mongo.Client
	Database string
}

// StartDatabase makes the connection with database and returns the connection
func StartDatabase() *DB {
	ctx := context.TODO()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("URL_URI")))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &DB{
		client:   c,
		Database: os.Getenv("DB_DATABASE"),
	}
}

// FormatUserCollection Formats the collection for users
func (db *DB) FormatUserCollection() *mongo.Collection {
	return db.client.Database(db.Database).Collection(os.Getenv("DB_USERS"))
}

// FormatGroupCollection Formats the collection for groupd
func (db *DB) FormatGroupCollection() *mongo.Collection {
	return db.client.Database(db.Database).Collection(os.Getenv("DB_GROUPD"))
}

// FormatUserChatlogs Formats the collection for user chatlogs
func (db *DB) FormatUserChatlogs() *mongo.Collection {
	return db.client.Database(db.Database).Collection(os.Getenv("DB_USR_CHLOGS"))
}

// FormatGroupChatlogs Formats the collection for group chat logs
func (db *DB) FormatGroupChatlogs() *mongo.Collection {
	return db.client.Database(db.Database).Collection(os.Getenv("DB_GR_CHLOGS"))
}
