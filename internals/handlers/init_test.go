package handlers

import (
	"wechat-back/internals/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MockDBName   = "test_db"
	MockObjectID = primitive.NewObjectID()
)

type DBMock struct {
	Client             *mongo.Client
	DatabaseName       string
	FindUserMockFunc   func(string) (models.User, bool, error)
	InsertUserMockFunc func(models.User) (string, error)
	UpdateUserMockFunc func(map[string]any, string) error
	GetUsersMockFunc   func(int, string) ([]*models.User, error)
}

func (db *DBMock) FindUserDB(email string) (models.User, bool, error) {

	if db.FindUserMockFunc != nil {
		return db.FindUserMockFunc(email)
	}
	return models.User{}, false, nil
}

func (db *DBMock) InsertUserDB(u models.User) (string, error) {
	if db.InsertUserMockFunc != nil {
		return db.InsertUserMockFunc(u)
	}
	return "", nil
}
func (db *DBMock) UpdateUserAccountDB(update map[string]any, id string) error {
	if db.UpdateUserMockFunc != nil {
		return db.UpdateUserMockFunc(update, id)
	}
	return nil
}
func (db *DBMock) GetUsers(pg int, query string) ([]*models.User, error) {
	if db.GetUsersMockFunc != nil {
		return db.GetUsersMockFunc(pg, query)
	}
	return []*models.User{}, nil
}
