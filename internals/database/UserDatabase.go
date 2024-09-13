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
FindUserDB
finds if a user has already an account
*/
func (db *DB) FindUserDB(email string) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"email": bson.M{"$eq": email},
	}

	var res models.User

	err := db.FormatUserCollection().FindOne(ctx, filter, nil).Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

/*
InsertUserDB
Will insert a new user to the collection
*/
func (db *DB) InsertUserDB(user models.User) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := db.FormatUserCollection().InsertOne(ctx, user, nil)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

/*
	UpdateUserAccountDB
	Updates the specified fileds that are passed as parameters on the user document
*/

func (db *DB) UpdateUserAccountDB(update map[string]any, i string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(i)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": bson.M{"$eq": id},
	}

	updateDoc := bson.M{
		"$set": update,
	}

	res, err := db.FormatUserCollection().UpdateOne(ctx, filter, updateDoc)
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
GetUsers
Will a list of users, the limit of each request is about 12
*/
func (db *DB) GetUsers(pg int, query string) ([]*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	filter := bson.M{
		"name": bson.M{
			"$regex": primitive.Regex{Pattern: query, Options: "i"},
		},
	}

	opts := options.Find()
	opts.SetLimit(12)
	opts.SetSkip(int64((pg - 1) * 12))

	var results []*models.User

	cursor, err := db.FormatUserCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {

		var user models.User

		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}

		results = append(results, &user)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}
