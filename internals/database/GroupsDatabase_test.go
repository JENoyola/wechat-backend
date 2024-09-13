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

func TestGetGroupDB(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetGroupDB - Success", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		i := ObjectIDMockHex

		group_name := "Funny Memes"

		group := bson.D{
			{Key: "name", Value: group_name},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.groups", mtest.FirstBatch, group))

		res, err := db.GetGroupDB(i)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.Name, group_name)
	})

	mt.Run("GetGroupDB  - Invalid object ID", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		i := "NOT A OBJECT ID"

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.groups", mtest.FirstBatch, nil))

		res, err := db.GetGroupDB(i)

		assert.EqualError(t, err, primitive.ErrInvalidHex.Error())
		assert.Nil(t, res)
	})

	mt.Run("GetGroupDB  - No documents", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		i := ObjectIDMockHex

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1235456,
			Message: mongo.ErrNoDocuments.Error(),
		}))

		res, err := db.GetGroupDB(i)

		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
		assert.Nil(t, res)
	})

}

// TestInsertGroupDB test database method InsertGroupDB
func TestInsertGroupDB(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("InsertGroupDB - Succes", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		group := models.Group{
			Name: "Funny memes",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		res, err := db.InsertGroupDB(group)

		assert.NoError(t, err)
		assert.NotEmpty(t, res)

	})

	mt.Run("InsertGroupDB - Error", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		group := models.Group{
			Name: "Funny memes",
		}

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    123456,
			Message: "Could not insert group",
		}))
		res, err := db.InsertGroupDB(group)

		assert.Error(t, err)
		assert.Empty(t, res)

	})

}

// TestTestUpdateGroupDB test database method TestUpdateGroupDB
func TestUpdateGroupDB(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("TestUpdateGroupDB - Succes", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbres := []bson.E{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 1},
		}

		update := map[string]any{
			"body": 5,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbres...))

		err := db.UpdateGroupDB(update, ObjectIDMock)

		assert.NoError(t, err)
	})

	mt.Run("TestUpdateGroupDB - Error no matched documents", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbres := []bson.E{
			{Key: "n", Value: 0},
			{Key: "nModified", Value: 1},
		}

		update := map[string]any{
			"body": 5,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbres...))

		err := db.UpdateGroupDB(update, ObjectIDMock)

		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
	})

	mt.Run("TestUpdateGroupDB - Error no modified", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbres := []bson.E{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 0},
		}

		update := map[string]any{
			"body": 5,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbres...))

		err := db.UpdateGroupDB(update, ObjectIDMock)

		assert.EqualError(t, err, ErrNoModified.Error())
	})

}

// TestDeleteGroupDB test database method DeleteGroupDB
func TestDeleteGroupDB(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("DeleteGroupDB - Success", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbRes := bson.E{
			Key:   "n",
			Value: 1,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbRes))
		err := db.DeleteGroupDB(ObjectIDMockHex)

		assert.NoError(t, err)

	})

	mt.Run("DeleteGroupDB - Error primitive", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbRes := bson.E{
			Key:   "n",
			Value: 1,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbRes))
		err := db.DeleteGroupDB("ijuhbsuihygs")

		assert.EqualError(t, err, primitive.ErrInvalidHex.Error())

	})

	mt.Run("DeleteGroupDB - Error no deleted", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbRes := bson.E{
			Key:   "n",
			Value: 0,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbRes))
		err := db.DeleteGroupDB(ObjectIDMockHex)

		assert.EqualError(t, err, ErrNoDeleted.Error())

	})

}

// TestSearchGroups test database method SearchGroups
func TestSearchGroups(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("SearchGroups - Success with Results", func(mt *mtest.T) {
		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		page := 1
		query := "group"

		groups := []bson.D{
			{{Key: "name", Value: "Group One"}, {Key: "description", Value: "First Group"}},
			{{Key: "name", Value: "Group Two"}, {Key: "description", Value: "Second Group"}},
		}

		// Mock a cursor response with two batches of documents.
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test_db.GROUPS", mtest.FirstBatch, groups...),
		)

		res, err := db.SearchGroups(page, query)

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "Group One", res[0].Name)
		assert.Equal(t, "Group Two", res[1].Name)
	})

	mt.Run("SearchGroups - No groups", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		pg := 1

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test_db.GROPS", mtest.FirstBatch, nil))

		res, err := db.SearchGroups(pg, "")

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	mt.Run("SearchGroups - encounter error", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		pg := 1

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    123456,
			Message: "Mongo db has encounter and error",
		}))

		res, err := db.SearchGroups(pg, "")

		assert.Error(t, err)
		assert.Nil(t, res)

	})

}
