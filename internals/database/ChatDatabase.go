package database

import (
	"context"
	"errors"
	"time"
	"wechat-back/internals/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	GetPrivateChatLogsDB
	Gets all the private message of a private converation
*/

func (db *DB) GetPrivateChatLogsDB(pg int, tar, curr string) ([]*models.P2PTextChatLog, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	current, err := primitive.ObjectIDFromHex(curr)
	if err != nil {
		return nil, err
	}

	target, err := primitive.ObjectIDFromHex(tar)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"$or": bson.A{
			bson.M{"target_id": bson.M{"$eq": target}, "author_id": bson.M{"$eq": current}},
			bson.M{"author_id": bson.M{"$eq": target}, "target_id": bson.M{"$eq": current}},
		},
	}

	opts := options.Find()
	opts.SetSkip(int64((pg - 1) * 60))
	opts.SetLimit(60)

	var res []*models.P2PTextChatLog

	cursor, err := db.FormatUserChatlogs().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {

		var chat models.P2PTextChatLog

		err := cursor.Decode(&chat)
		if err != nil {
			return nil, err
		}

		res = append(res, &chat)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

/*
GetGroupChatLogsDB
gets the chat of a specific group
*/
func (db *DB) GetGroupChatLogsDB(pg int, groupid string) ([]*models.GroupChatTextLog, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(groupid)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": bson.M{"$eq": id},
	}

	opts := options.Find()
	opts.SetSkip(int64((pg - 1) * 60))
	opts.SetLimit(60)

	var res []*models.GroupChatTextLog

	cursor, err := db.FormatGroupChatlogs().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {

		var chat models.GroupChatTextLog

		err := cursor.Decode(&chat)
		if err != nil {
			return nil, err
		}

		res = append(res, &chat)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	return res, nil

}

/*
	InsertP2PMessageDB
	Inserts a new message to the private conversation collection
*/

func (db *DB) InsertP2PMessageDB(m models.P2PTextChatLog) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	info, err := db.FormatUserChatlogs().InsertOne(ctx, m, nil)
	if err != nil {
		return "", err
	}

	return info.InsertedID.(primitive.ObjectID).Hex(), nil

}

/*
InsertGroupMessageDB
Inserts a new message to the group collection
*/
func (db *DB) InsertGroupMessageDB(m models.GroupChatTextLog) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := db.FormatGroupChatlogs().InsertOne(ctx, m, nil)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

/*
UpdateP2PMessageDB
updates the given P2P message on the database
*/
func (db *DB) UpdateP2PMessageDB(updated map[string]any, chtid string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(chtid)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": bson.M{"$eq": id},
	}

	updateDoc := bson.M{
		"$set": updated,
	}

	res, err := db.FormatUserChatlogs().UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	if res.MatchedCount < 1 {
		return mongo.ErrNoDocuments
	} else if res.ModifiedCount < 1 {
		return errors.New("no documents modified")
	}

	return nil
}

/*
	UpdateGroupMessageDB
	updates the given group message on the database
*/

func (db *DB) UpdateGroupMessageDB(updated map[string]any, chtid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(chtid)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": bson.M{"$eq": id},
	}

	updateDoc := bson.M{
		"$set": updated,
	}

	res, err := db.FormatGroupChatlogs().UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	} else if res.MatchedCount < 1 {
		return mongo.ErrNoDocuments
	} else if res.ModifiedCount < 1 {
		return errors.New("no documents modified")
	}

	return nil

}
