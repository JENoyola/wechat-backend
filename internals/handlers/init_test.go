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

	// groups
	GetGroupDBMockFunc    func(string) (*models.Group, error)
	InsertGroupDBMockFunc func(models.Group) (string, error)
	UpdateGroupDBMockFunc func(map[string]any, primitive.ObjectID) error
	DeleteGroupDBMockFunc func(string) error
	SearchGroupsMockFunc  func(int, string) ([]*models.Group, error)

	// Chat
	InsertP2PMessageDBMockFunc  func(any) (string, error)
	InsertGroupMessageDBMockFun func(any) (string, error)
}

/*USER MOCK FUNCTIONS*/

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

/*GROUP MOCK FUNCTIONS*/
func (db *DBMock) GetGroupDB(s string) (*models.Group, error) {
	if db.GetGroupDBMockFunc != nil {
		return db.GetGroupDBMockFunc(s)
	}
	return &models.Group{}, nil
}

func (db *DBMock) InsertGroupDB(g models.Group) (string, error) {
	if db.InsertGroupDBMockFunc != nil {
		return db.InsertGroupDBMockFunc(g)
	}
	return "", nil
}

func (db *DBMock) UpdateGroupDB(u map[string]any, i primitive.ObjectID) error {
	if db.UpdateGroupDBMockFunc != nil {
		return db.UpdateGroupDBMockFunc(u, i)
	}
	return nil
}

func (db *DBMock) DeleteGroupDB(i string) error {
	if db.DeleteGroupDBMockFunc != nil {
		return db.DeleteGroupDBMockFunc(i)
	}
	return nil
}

func (db *DBMock) SearchGroups(pg int, query string) ([]*models.Group, error) {
	if db.SearchGroupsMockFunc != nil {
		return db.SearchGroupsMockFunc(pg, query)
	}
	return []*models.Group{}, nil
}

// CHAT METHODS

func (db *DBMock) InsertP2PMessageDB(m any) (string, error) {
	if db.InsertP2PMessageDBMockFunc != nil {
		return db.InsertP2PMessageDBMockFunc(m)
	}
	return "", nil
}

func (db *DBMock) InsertGroupMessageDB(m any) (string, error) {
	if db.InsertGroupMessageDBMockFun != nil {
		return db.InsertGroupMessageDBMockFun(m)
	}
	return "", nil
}
