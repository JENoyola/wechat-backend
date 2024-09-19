package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wechat-back/internals/decorators"
	"wechat-back/internals/models"
	"wechat-back/internals/server"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TESTS TO BE UPLOADED
func TestNewUserAccountEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("NewUserAccountEP  - Success insertion", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"

		user := &models.User{
			Name:  "jorge",
			Email: expectedEmail,
		}

		db := &DBMock{
			DatabaseName: MockDBName,
			Client:       mt.Client,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return models.User{}, false, nil
			},
			InsertUserMockFunc: func(u models.User) (string, error) {
				return "", nil
			},
		}
		bod, err := json.Marshal(user)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/nsg", bytes.NewReader(bod))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(NewUserAccountEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectedEmail, res.DATA)

	})

	mt.Run("NewUserAccountEP - Error user exist", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"

		user := &models.User{
			Name:  "jorge",
			Email: expectedEmail,
		}

		db := &DBMock{
			DatabaseName: MockDBName,
			Client:       mt.Client,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return models.User{}, true, nil
			},
			InsertUserMockFunc: func(u models.User) (string, error) {
				return "", nil
			},
		}
		bod, err := json.Marshal(user)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/nsg", bytes.NewReader(bod))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(NewUserAccountEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	mt.Run("NewUserAccountEP - DB Error", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"

		user := &models.User{
			Name:  "jorge",
			Email: expectedEmail,
		}

		db := &DBMock{
			DatabaseName: MockDBName,
			Client:       mt.Client,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return models.User{}, false, errors.New("could not find user")
			},
			InsertUserMockFunc: func(u models.User) (string, error) {
				return "", nil
			},
		}
		bod, err := json.Marshal(user)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/nsg", bytes.NewReader(bod))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(NewUserAccountEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusConflict, rr.Code)

	})

	mt.Run("NewUserAccountEP - failed insertion", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"

		user := &models.User{
			Name:  "jorge",
			Email: expectedEmail,
		}

		db := &DBMock{
			DatabaseName: MockDBName,
			Client:       mt.Client,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return models.User{}, false, nil
			},
			InsertUserMockFunc: func(u models.User) (string, error) {
				return "", errors.New("could not insert document")
			},
		}
		bod, err := json.Marshal(user)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/nsg", bytes.NewReader(bod))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(NewUserAccountEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
	})

}

func TestRequestNewCodeEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("RequestNewCodeEP - Success", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedUser := &models.User{
			ID:    MockObjectID,
			Name:  "jorge",
			Email: expectedEmail,
		}

		// setup
		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, nil
			},
			UpdateUserMockFunc: func(m map[string]any, s string) error {
				return nil
			},
		}

		// Set up HTTP request for handler
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sn?e=%s", expectedEmail), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(RequestNewCodeEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, server.OK, res.Code)
	})

	mt.Run("RequestNewCodeEP - Error user not find ", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedUser := &models.User{
			ID:    MockObjectID,
			Name:  "jorge",
			Email: expectedEmail,
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, false, nil
			},
			UpdateUserMockFunc: func(m map[string]any, s string) error {
				return nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sn?e=%s", expectedEmail), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(RequestNewCodeEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, server.NO_DOCUMENTS, res.Code)
	})

	mt.Run("RequestNewCodeEP - Error database find ", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedUser := &models.User{
			ID:    MockObjectID,
			Name:  "jorge",
			Email: expectedEmail,
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, errors.New("could not complete operation")
			},
			UpdateUserMockFunc: func(m map[string]any, s string) error {
				return nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sn?e=%s", expectedEmail), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(RequestNewCodeEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, server.DB_ERROR, res.Code)
	})

	mt.Run("RequestNewCodeEP - Error updating user document ", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedUser := &models.User{
			ID:    MockObjectID,
			Name:  "jorge",
			Email: expectedEmail,
		}

		expectedErr := errors.New("error updating database")

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, nil
			},
			UpdateUserMockFunc: func(m map[string]any, s string) error {
				return expectedErr
			},
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sn?e=%s", expectedEmail), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(RequestNewCodeEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.EqualError(t, expectedErr, res.Message)
	})

}

func TestUserCodeVerificationEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("UserCodeVerificationEP - Success verification", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedCode := 123456

		expectedUser := &models.User{
			ID:            MockObjectID,
			Email:         expectedEmail,
			TempCode:      expectedCode,
			ValidCode:     models.CODE_VALID,
			CodeTimestamp: time.Now(),
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, nil
			},
			UpdateUserMockFunc: func(m map[string]any, s string) error {
				return nil
			},
		}

		payload := &models.User{
			Email:    expectedEmail,
			TempCode: expectedCode,
		}

		bod, err := json.Marshal(payload)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, "/cv", bytes.NewReader(bod))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UserCodeVerificationEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.NotNil(t, res.DATA)

	})

	mt.Run("UserCodeVerificationEP - Error findOne", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedCode := 123456
		expectedError := errors.New("something went wrong finding document")

		expectedUser := &models.User{
			ID:            MockObjectID,
			Email:         expectedEmail,
			TempCode:      expectedCode,
			ValidCode:     models.CODE_VALID,
			CodeTimestamp: time.Now(),
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, expectedError
			},
		}

		payload := &models.User{
			Email:    expectedEmail,
			TempCode: expectedCode,
		}

		bod, err := json.Marshal(payload)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, "/cv", bytes.NewReader(bod))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UserCodeVerificationEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedError, res.Message)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UserCodeVerificationEP - User not exist", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedCode := 123456
		expectedError := errors.New("not allowed")

		expectedUser := &models.User{
			ID:            MockObjectID,
			Email:         expectedEmail,
			TempCode:      expectedCode,
			ValidCode:     models.CODE_VALID,
			CodeTimestamp: time.Now(),
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, false, nil
			},
		}

		payload := &models.User{
			Email:    expectedEmail,
			TempCode: expectedCode,
		}

		bod, err := json.Marshal(payload)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, "/cv", bytes.NewReader(bod))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UserCodeVerificationEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedError, res.Message)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UserCodeVerificationEP - Error Fail update", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedCode := 123456
		expectedError := errors.New("could not update user")

		expectedUser := &models.User{
			ID:            MockObjectID,
			Email:         expectedEmail,
			TempCode:      expectedCode,
			ValidCode:     models.CODE_VALID,
			CodeTimestamp: time.Now(),
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, nil
			},
			UpdateUserMockFunc: func(m map[string]any, s string) error {
				return expectedError
			},
		}

		payload := &models.User{
			Email:    expectedEmail,
			TempCode: expectedCode,
		}

		bod, err := json.Marshal(payload)
		assert.Nil(t, err)

		req := httptest.NewRequest(http.MethodPut, "/cv", bytes.NewReader(bod))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UserCodeVerificationEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedError, res.Message)
		assert.Nil(t, res.DATA)

	})

}

func TestLoginEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("LoginEP - Success login ", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedUser := &models.User{
			ID:    MockObjectID,
			Email: expectedEmail,
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, nil
			},
			UpdateUserMockFunc: func(m map[string]any, s string) error {
				return nil
			},
		}

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/l?e=%s", expectedEmail), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(LoginEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.REDIRECTION, res.Code)
		assert.NotNil(t, res.DATA)
	})

	// error database findone
	mt.Run("LoginEP - Error findOne ", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedUser := &models.User{}
		expectedError := errors.New("error retrieving document")

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, expectedError
			},
		}

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/l?e=%s", expectedEmail), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(LoginEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_FIELD, res.Code)
		assert.Nil(t, res.DATA)
	})
	// error user not exist
	mt.Run("LoginEP - user not exist ", func(mt *mtest.T) {

		expectedEmail := ""
		expectedUser := &models.User{}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, false, nil
			},
		}

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/l?e=%s", expectedEmail), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(LoginEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.NOT_ALLOWED, res.Code)
		assert.Nil(t, res.DATA)
	})
	// error updating
	mt.Run("LoginEP - Error updating ", func(mt *mtest.T) {

		expectedEmail := "jorge@mail.com"
		expectedUser := &models.User{
			ID:    MockObjectID,
			Email: expectedEmail,
		}
		expectedError := errors.New("error updating document")

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return *expectedUser, true, nil
			},
			UpdateUserMockFunc: func(m map[string]any, s string) error {
				return expectedError
			},
		}

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/l?e=%s", expectedEmail), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(LoginEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)
	})
}

func TestSearchUsersEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("SearchUsersEP - Success search", func(mt *mtest.T) {

		expectedResults := []*models.User{}

		for i := 1; i <= 12; i++ {
			usr := &models.User{
				ID:    MockObjectID,
				Name:  "Jorge",
				Email: "jorge@mail.com",
			}

			expectedResults = append(expectedResults, usr)
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetUsersMockFunc: func(i int, s string) ([]*models.User, error) {
				return expectedResults, nil
			},
		}

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/l?pg=%d&q=%s", 1, ""), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(SearchUsersEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.OK, res.Code)
		assert.NotNil(t, res.DATA)
		assert.Len(t, expectedResults, 12)
	})

	// No users
	mt.Run("SearchUsersEP - no users", func(mt *mtest.T) {

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetUsersMockFunc: func(i int, s string) ([]*models.User, error) {
				return nil, nil
			},
		}

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/l?pg=%d&q=%s", 1, ""), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(SearchUsersEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.OK, res.Code)
		assert.Nil(t, res.DATA)
	})

	// db error
	mt.Run("SearchUsersEP - database error", func(mt *mtest.T) {

		expectedError := errors.New("could not retrive results")

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetUsersMockFunc: func(i int, s string) ([]*models.User, error) {
				return nil, expectedError
			},
		}

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/l?pg=%d&q=%s", 1, ""), nil)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(SearchUsersEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err := json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)
	})
}
