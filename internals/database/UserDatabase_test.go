package database

import (
	"testing"
	"wechat-back/internals/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestFindUserDB tests the FindUserDB method
func TestFindUserDB(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("FindUser - Success", func(mt *mtest.T) {
		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		email := "jorge@mail.com"
		user := models.User{Email: email}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test_db.users", mtest.FirstBatch, bson.D{
			{Key: "email", Value: user.Email},
		}))

		res, err := db.FindUserDB(email)
		assert.NoError(t, err)
		assert.Equal(t, email, res.Email)
	})

	mt.Run("FindUser - Not Found", func(mt *mtest.T) {
		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		phone := "jorge@mail.com"
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.users", mtest.FirstBatch))

		_, err := db.FindUserDB(phone)
		assert.Error(t, err)
	})
}

// TestInsertUserDB tests the InsertUserDB method
func TestInsertUserDB(t *testing.T) {

	mongoTest := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mongoTest.Run("InsertUser - Success", func(mt *mtest.T) {
		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		user := models.User{
			Name:      "George",
			Email:     "jorge@mail.com",
			ValidCode: 0,
			TempCode:  123456,
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		id, err := db.InsertUserDB(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

	})

	mongoTest.Run("InsertUser - Error", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		user := models.User{
			Name:      "George",
			Email:     "jorge@mail.com",
			ValidCode: 0,
			TempCode:  123456,
		}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Message: "Duplicate user",
		}))
		_, err := db.InsertUserDB(user)
		assert.Error(t, err)
	})

}

// TestUpdateUserAccountDB tests the UpdateUserAccountDB method
func TestUpdateUserAccountDB(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("UpdateUser - Success", func(mt *mtest.T) {
		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbRes := []bson.E{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 1},
		}

		update := map[string]interface{}{
			"first_name": "John",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbRes...))
		err := db.UpdateUserAccountDB(update, ObjectIDMockHex)
		assert.NoError(t, err)
	})

	mt.Run("UpdateUser - No Match", func(mt *mtest.T) {
		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbRes := []bson.E{
			{Key: "n", Value: 0},
			{Key: "nModified", Value: 1},
		}

		update := map[string]interface{}{
			"first_name": "John",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbRes...))
		err := db.UpdateUserAccountDB(update, ObjectIDMockHex)
		assert.EqualError(t, err, mongo.ErrNoDocuments.Error())
	})

	mt.Run("UpdateUser - No Match", func(mt *mtest.T) {
		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		dbRes := []bson.E{
			{Key: "n", Value: 1},
			{Key: "nModified", Value: 0},
		}

		update := map[string]interface{}{
			"first_name": "John",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse(dbRes...))
		err := db.UpdateUserAccountDB(update, ObjectIDMockHex)
		assert.EqualError(t, err, ErrNoModified.Error())
	})

}

// TestGetUsers test the GetUser methods
func TestGetUsers(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetUsers - Success", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		pg := 1

		users := []bson.D{
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "John"}, {Key: "last_name", Value: "Doe"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
			{{Key: "email", Value: "jorge@mail.com"}, {Key: "first_name", Value: "Jane"}, {Key: "last_name", Value: "Smith"}},
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test_db.users", mtest.FirstBatch, users...))

		res, err := db.GetUsers(pg, "")

		assert.NoError(t, err)
		assert.Len(t, res, 12)
		assert.Equal(t, "jorge@mail.com", res[0].Email)
	})

	mt.Run("GetUsers - No users", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		pg := 1

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test_db.users", mtest.FirstBatch, nil))

		res, err := db.GetUsers(pg, "")

		assert.Error(t, err)
		assert.Empty(t, res)
	})

	mt.Run("GetUsers - encounter error", func(mt *mtest.T) {

		db := &DB{
			client:   mt.Client,
			Database: mockDBName,
		}

		pg := 1

		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    123456,
			Message: "Mongo db has encounter and error",
		}))

		res, err := db.GetUsers(pg, "")

		assert.Error(t, err)
		assert.Nil(t, res)

	})

}
