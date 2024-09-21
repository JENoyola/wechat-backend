package database

import (
	"context"
	"time"
	"wechat-back/internals/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
GetGroupDB
Gets a group by the ID
*/
func (db *DB) GetGroupDB(i string) (*models.Group, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"group_id": bson.M{"$eq": i},
	}

	var res models.Group

	err := db.FormatGroupCollection().FindOne(ctx, filter, nil).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

/*
InsertGroupDB
Inserts a new group document to the database
*/
func (db *DB) InsertGroupDB(g models.Group) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := db.FormatGroupCollection().InsertOne(ctx, g, nil)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

/*
UpdateGroupDB
Updates the document passing the updated fields
*/
func (db *DB) UpdateGroupDB(update map[string]any, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": bson.M{"$eq": id},
	}

	updatedDoc := bson.M{
		"$set": update,
	}

	res, err := db.FormatGroupCollection().UpdateOne(ctx, filter, updatedDoc)
	if err != nil {
		return err
	}

	if res.MatchedCount < 1 {
		return mongo.ErrNoDocuments
	} else if res.ModifiedCount < 1 {
		return ErrNoModified
	}

	return nil
}

/*
	DeleteGroupDB
	Deletes a group document on the collection selected
*/

func (db *DB) DeleteGroupDB(i string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(i)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": bson.M{"$eq": id},
	}

	res, err := db.FormatGroupCollection().DeleteOne(ctx, filter, nil)
	if err != nil {
		return err
	}

	if res.DeletedCount < 1 {
		return ErrNoDeleted
	}

	return nil
}

/*
	SearchGroups
	Gets a list of groups matching the query or deliberately
*/

func (db *DB) SearchGroups(page int, query string) ([]*models.Group, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	filter := bson.M{
		"name": bson.M{
			"$regex": primitive.Regex{Pattern: query, Options: "i"},
		},
	}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * 8))
	opts.SetLimit(8)

	var res []*models.Group

	cursor, err := db.FormatGroupCollection().Find(ctx, filter, opts)
	if err != nil {
		return res, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {

		var group models.Group

		err := cursor.Decode(&group)
		if err != nil {
			return res, err
		}

		res = append(res, &group)
	}

	err = cursor.Err()
	if err != nil {
		return res, err
	}

	return res, nil
}
