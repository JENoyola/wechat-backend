package database

import (
	"testing"
	"wechat-back/internals/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestGetPrivateChatLogsDB test database method GetPrivateChatLogsDB
func TestGetPrivateChatLogsDB(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetPrivateChatLogsDB - Success", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		pg := 1
		tar := "66d6561e43416dd7f7eb6aa5"
		curr := ObjectIDMockHex

		Logs := []bson.D{
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
			{{Key: "target_id", Value: curr}, {Key: "author_id", Value: tar}, {Key: "body", Value: "Hola, tanto sin escuchar de ti"}},
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.USCHLOGS", mtest.FirstBatch, Logs...))

		res, err := db.GetPrivateChatLogsDB(pg, tar, curr)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res, 4)

	})

	mt.Run("GetPrivateChatLogsDB - primitive error", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		pg := 1
		tar := "Not a primitive id"
		curr := ObjectIDMockHex

		Logs := []bson.D{
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
			{{Key: "target_id", Value: curr}, {Key: "author_id", Value: tar}, {Key: "body", Value: "Hola, tanto sin escuchar de ti"}},
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.USCHLOGS", mtest.FirstBatch, Logs...))

		res, err := db.GetPrivateChatLogsDB(pg, tar, curr)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Len(t, res, 0)

	})

	mt.Run("GetPrivateChatLogsDB - No documents", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		pg := 1
		tar := "66d6561e43416dd7f7eb6aa5"
		curr := ObjectIDMockHex

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.USCHLOGS", mtest.FirstBatch, nil))

		res, err := db.GetPrivateChatLogsDB(pg, tar, curr)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Len(t, res, 0)

	})

}

// TestGetGroupChatLogsDB test database method GetGroupChatLogsDB
func TestGetGroupChatLogsDB(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetGroupChatLogsDB - Success", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		pg := 1
		tar := "66d6561e43416dd7f7eb6aa5"
		curr := ObjectIDMockHex

		Logs := []bson.D{
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
			{{Key: "target_id", Value: curr}, {Key: "author_id", Value: tar}, {Key: "body", Value: "Hola, tanto sin escuchar de ti"}},
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.USCHLOGS", mtest.FirstBatch, Logs...))

		res, err := db.GetGroupChatLogsDB(pg, curr)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res, 4)

	})

	mt.Run("GetGroupChatLogsDB - primitive error", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		pg := 1
		tar := "Not a primitive id"
		curr := ObjectIDMockHex

		Logs := []bson.D{
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
			{{Key: "target_id", Value: curr}, {Key: "author_id", Value: tar}, {Key: "body", Value: "Hola, tanto sin escuchar de ti"}},
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
			{{Key: "target_id", Value: tar}, {Key: "author_id", Value: curr}, {Key: "body", Value: "Hola, como estas"}},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.USCHLOGS", mtest.FirstBatch, Logs...))

		res, err := db.GetGroupChatLogsDB(pg, curr)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Len(t, res, 0)

	})

	mt.Run("GetGroupChatLogsDB - No documents", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		pg := 1
		curr := ObjectIDMockHex

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.USCHLOGS", mtest.FirstBatch, nil))

		res, err := db.GetGroupChatLogsDB(pg, curr)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Len(t, res, 0)

	})

}

// TestInsertP2PMessageDB test database method InsertP2PMessageDB
func TestInsertP2PMessageDB(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("InsertP2PMessageDB - Success", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		m := models.P2PChatLog{
			Body: "Hola como estas",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		id, err := db.InsertP2PMessageDB(m)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)

	})

	mt.Run("InsertP2PMessageDB - Error", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		m := models.P2PChatLog{
			Body: "Hola como estas",
		}

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    12345,
			Message: "Error inserting message",
		}))
		id, err := db.InsertP2PMessageDB(m)

		assert.Error(t, err)
		assert.Empty(t, id)

	})

}

// TestInsertGroupMessageDB test database method InsertGroupMessageDB
func TestInsertGroupMessageDB(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("InsertGroupMessageDB - Success", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		m := models.GroupChatLog{
			Body: "Hola como estas",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		id, err := db.InsertGroupMessageDB(m)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)

	})

	mt.Run("InsertGroupMessageDB - Error", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		m := models.GroupChatLog{
			Body: "Hola como estas",
		}

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    12345,
			Message: "Error inserting message",
		}))
		id, err := db.InsertGroupMessageDB(m)

		assert.Error(t, err)
		assert.Empty(t, id)

	})

}

// TestUpdateP2PMessageDB test database method UpdateP2PMessageDB
func TestUpdateP2PMessageDB(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("UpdateP2PMessage - Success", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		res := []bson.E{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 1},
		}

		chtid := "66d6561e43416dd7f7eb6aa5"
		updated := map[string]interface{}{
			"body": "Updated message",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(res...))

		err := db.UpdateP2PMessageDB(updated, chtid)

		assert.NoError(t, err) // Expect no error

	})

	mt.Run("UpdateP2PMessage - BSON error", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		update := make(map[string]any)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := db.UpdateP2PMessageDB(update, "skjs")

		assert.EqualError(t, err, primitive.ErrInvalidHex.Error())
	})

	mt.Run("UpdateP2PMessage - No documents", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		res := []bson.E{
			{Key: "n", Value: 0},
			{Key: "nModified", Value: 1},
		}

		update := make(map[string]any)

		mt.AddMockResponses(mtest.CreateSuccessResponse(res...))
		err := db.UpdateP2PMessageDB(update, ObjectIDMockHex)
		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
	})

	mt.Run("UpdateP2PMessage - Error no modified", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		res := []bson.E{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 0},
		}

		update := make(map[string]any)

		mt.AddMockResponses(mtest.CreateSuccessResponse(res...))
		err := db.UpdateP2PMessageDB(update, ObjectIDMockHex)
		assert.EqualError(t, err, "no documents modified")
	})

}

// TestUpdateGroupMessageDB test database method UpdateGroupMessageDB
func TestUpdateGroupMessageDB(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("UpdateGroupMessageDB - Success", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		res := []bson.E{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 1},
		}

		chtid := "66d6561e43416dd7f7eb6aa5"
		updated := map[string]interface{}{
			"body": "Updated message",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(res...))

		err := db.UpdateGroupMessageDB(updated, chtid)

		assert.NoError(t, err) // Expect no error

	})

	mt.Run("UpdateGroupMessageDB - BSON error", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		update := make(map[string]any)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := db.UpdateGroupMessageDB(update, "skjs")

		assert.EqualError(t, err, primitive.ErrInvalidHex.Error())
	})

	mt.Run("UpdateGroupMessageDB - No documents", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		res := []bson.E{
			{Key: "n", Value: 0},
			{Key: "nModified", Value: 1},
		}

		update := make(map[string]any)

		mt.AddMockResponses(mtest.CreateSuccessResponse(res...))
		err := db.UpdateGroupMessageDB(update, ObjectIDMockHex)
		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
	})

	mt.Run("UpdateGroupMessageDB - Error no modified", func(mt *mtest.T) {

		db := &DB{
			Client:   mt.Client,
			Database: MockDBName,
		}

		res := []bson.E{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 0},
		}

		update := make(map[string]any)

		mt.AddMockResponses(mtest.CreateSuccessResponse(res...))
		err := db.UpdateGroupMessageDB(update, ObjectIDMockHex)
		assert.EqualError(t, err, "no documents modified")
	})

}
