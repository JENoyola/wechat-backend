package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wechat-back/internals/decorators"
	"wechat-back/internals/models"
	"wechat-back/internals/server"
	"wechat-back/providers/media"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestHandleP2PConnectionEP test the handler HandleP2PConnectionEP
func TestHandleP2PConnectionEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("HandleP2PConnectionEP - Successful connection and message sent", func(mt *mtest.T) {

		server.StartWebsocketService()

		expectedEmail := "george@mail.com"
		tar := primitive.NewObjectID()
		tar_token := "111111111"

		body := models.InboundP2PTextMessage{
			MessageID:       "",
			AuthorID:        MockObjectID.Hex(),
			TargetID:        tar.Hex(),
			TargetPushToken: tar_token,
			Body:            "Heyy how you are doing",
		}

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, true, nil
			},
			InsertP2PMessageDBMockFunc: func(ppcl any) (string, error) {

				return "", nil
			},
		}

		m := &media.MediaMock{}

		server := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleP2PConnectionEP, db, m)))
		defer server.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?&email=%s&tar=%s&tar_tkn=%s", strings.ReplaceAll(server.URL, "http", "ws"), expectedEmail, tar.Hex(), tar_token), nil)
		assert.Nil(t, err)
		defer conn.Close()

		data2Send, err := json.Marshal(body)
		assert.Nil(t, err)

		err = conn.WriteMessage(websocket.TextMessage, data2Send)
		assert.Nil(t, err)

		_, data, err := conn.ReadMessage()
		assert.Nil(t, err)

		var res models.OutboundP2PTextMessage

		err = json.Unmarshal(data, &res)
		assert.Nil(t, err)

		assert.Equal(t, body.Body, res.Body)

	})

	mt.Run("HandleP2PConnectionEP - Error author find one", func(mt *mtest.T) {

		server.StartWebsocketService()

		expectedError := errors.New("error finding user")

		expectedEmail := "george@mail.com"
		tar := primitive.NewObjectID()
		tar_token := "111111111"

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, false, expectedError
			},
			InsertP2PMessageDBMockFunc: func(ppcl any) (string, error) {

				return "", nil
			},
		}

		m := &media.MediaMock{}

		S := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleP2PConnectionEP, db, m)))
		defer S.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?&email=%s&tar=%s&tar_tkn=%s", strings.ReplaceAll(S.URL, "http", "ws"), expectedEmail, tar.Hex(), tar_token), nil)
		assert.Nil(t, err)
		defer conn.Close()

		var res models.OutboundP2PTextMessage

		err = conn.ReadJSON(&res)
		assert.Nil(t, err)

		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_REQUEST, res.Code)
		assert.EqualError(t, expectedError, res.Message)

	})

	mt.Run("HandleP2PConnectionEP - Error author not exist", func(mt *mtest.T) {

		server.StartWebsocketService()

		expectedEmail := "george@mail.com"
		tar := primitive.NewObjectID()
		tar_token := "111111111"

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, false, nil
			},
			InsertP2PMessageDBMockFunc: func(ppcl any) (string, error) {

				return "", nil
			},
		}

		m := &media.MediaMock{}

		S := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleP2PConnectionEP, db, m)))
		defer S.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?&email=%s&tar=%s&tar_tkn=%s", strings.ReplaceAll(S.URL, "http", "ws"), expectedEmail, tar.Hex(), tar_token), nil)
		assert.Nil(t, err)
		defer conn.Close()

		var res models.OutboundP2PTextMessage

		err = conn.ReadJSON(&res)
		assert.Nil(t, err)

		assert.True(t, res.Error)
		assert.Equal(t, server.NOT_ALLOWED, res.Code)

	})

	mt.Run("HandleP2PConnectionEP - Error json Inbound message", func(mt *mtest.T) {

		server.StartWebsocketService()

		expectedEmail := "george@mail.com"
		tar := primitive.NewObjectID()
		tar_token := "111111111"

		body := "Not a json Object"

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, true, nil
			},
			InsertP2PMessageDBMockFunc: func(ppcl any) (string, error) {

				return "", nil
			},
		}

		m := &media.MediaMock{}

		S := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleP2PConnectionEP, db, m)))
		defer S.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?&email=%s&tar=%s&tar_tkn=%s", strings.ReplaceAll(S.URL, "http", "ws"), expectedEmail, tar.Hex(), tar_token), nil)
		assert.Nil(t, err)
		defer conn.Close()

		err = conn.WriteMessage(websocket.TextMessage, []byte(body))
		assert.Nil(t, err)

		_, data, err := conn.ReadMessage()
		assert.Nil(t, err)

		var res models.OutboundP2PTextMessage

		err = json.Unmarshal(data, &res)
		assert.Nil(t, err)

		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_REQUEST, res.Code)

	})

	mt.Run("HandleP2PConnectionEP - Error invalid targetID", func(mt *mtest.T) {

		server.StartWebsocketService()

		expectedEmail := "george@mail.com"
		tar := "notAObjectID"
		tar_token := "111111111"

		body := models.InboundP2PTextMessage{
			MessageID:       "",
			AuthorID:        MockObjectID.Hex(),
			TargetID:        tar,
			TargetPushToken: tar_token,
			Body:            "Heyy how you are doing",
		}

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, true, nil
			},
			InsertP2PMessageDBMockFunc: func(ppcl any) (string, error) {

				return "", nil
			},
		}

		m := &media.MediaMock{}

		S := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleP2PConnectionEP, db, m)))
		defer S.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?&email=%s&tar=%s&tar_tkn=%s", strings.ReplaceAll(S.URL, "http", "ws"), expectedEmail, tar, tar_token), nil)
		assert.Nil(t, err)
		defer conn.Close()

		data2Send, err := json.Marshal(body)
		assert.Nil(t, err)

		err = conn.WriteMessage(websocket.TextMessage, data2Send)
		assert.Nil(t, err)

		_, data, err := conn.ReadMessage()
		assert.Nil(t, err)

		var res models.OutboundP2PTextMessage

		err = json.Unmarshal(data, &res)
		assert.Nil(t, err)

		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_FIELD, res.Code)
		assert.EqualError(t, primitive.ErrInvalidHex, res.Message)
	})
}

// TestHandleGroupConnectionsEP test the handler HandleGroupConnectionsEP
func TestHandleGroupConnectionsEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("HandleGroupConnectionsEP - Successful connection and message sent", func(mt *mtest.T) {

		server.StartWebsocketService()

		expectedEmail := "jorge@mail.com"
		expectedGroupID := "6177226702-5T2de426p8arbt6sb4b128o63afaG9u3f-1727206726"
		expectedUsers := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID(), MockObjectID}

		messageBody := models.InboundGroupTextMessage{
			MessageID: "",
			AuthorID:  MockObjectID.Hex(),
			GroupID:   expectedGroupID,
			PushTokens: map[string]string{
				expectedUsers[0].Hex(): "333333333",
				expectedUsers[1].Hex(): "222222222",
				expectedUsers[2].Hex(): "111111111",
			},
			Body: "Hey everyone",
		}

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		returnedGroup := models.Group{
			ID:           primitive.NewObjectID(),
			GroupID:      expectedGroupID,
			Name:         "Some group name",
			Description:  "Test group",
			Participants: expectedUsers,
			Admins:       []primitive.ObjectID{MockObjectID},
			ProfileImage: "",
		}

		db := DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, true, nil
			},
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return &returnedGroup, nil
			},
			InsertGroupMessageDBMockFun: func(m any) (string, error) {
				return "", nil
			},
		}

		m := &media.MediaMock{}

		S := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleGroupConnectionsEP, &db, m)))
		defer S.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?gi=%s&ui=%s", strings.ReplaceAll(S.URL, "http", "ws"), expectedGroupID, expectedEmail), nil)
		assert.Nil(t, err)
		defer conn.Close()

		data2Send, err := json.Marshal(messageBody)
		assert.Nil(t, err)

		err = conn.WriteMessage(websocket.TextMessage, data2Send)
		assert.Nil(t, err)

		_, data, err := conn.ReadMessage()
		assert.Nil(t, err)

		var res models.OutboundGroupTextMessage

		err = json.Unmarshal(data, &res)
		assert.Nil(t, err)

		assert.False(t, res.Error)
		assert.Equal(t, messageBody.Body, res.Body)

	})

	mt.Run("HandleGroupConnectionsEP - Error Not author found", func(mt *mtest.T) {
		server.StartWebsocketService()

		expectedError := mongo.ErrNoDocuments

		expectedEmail := "jorge@mail.com"
		expectedGroupID := "6177226702-5T2de426p8arbt6sb4b128o63afaG9u3f-1727206726"
		expectedUsers := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID(), MockObjectID}

		messageBody := models.InboundGroupTextMessage{
			MessageID: "",
			AuthorID:  MockObjectID.Hex(),
			GroupID:   expectedGroupID,
			PushTokens: map[string]string{
				expectedUsers[0].Hex(): "333333333",
				expectedUsers[1].Hex(): "222222222",
				expectedUsers[2].Hex(): "111111111",
			},
			Body: "Hey everyone",
		}

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		returnedGroup := models.Group{
			ID:           primitive.NewObjectID(),
			GroupID:      expectedGroupID,
			Name:         "Some group name",
			Description:  "Test group",
			Participants: expectedUsers,
			Admins:       []primitive.ObjectID{MockObjectID},
			ProfileImage: "",
		}

		db := DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, false, expectedError
			},
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return &returnedGroup, nil
			},
			InsertGroupMessageDBMockFun: func(m any) (string, error) {
				return "", nil
			},
		}

		m := &media.MediaMock{}

		S := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleGroupConnectionsEP, &db, m)))
		defer S.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?gi=%s&ui=%s", strings.ReplaceAll(S.URL, "http", "ws"), expectedGroupID, expectedEmail), nil)
		assert.Nil(t, err)
		defer conn.Close()

		data2Send, err := json.Marshal(messageBody)
		assert.Nil(t, err)

		err = conn.WriteMessage(websocket.TextMessage, data2Send)
		assert.Nil(t, err)

		_, data, err := conn.ReadMessage()
		assert.Nil(t, err)

		var res models.OutboundGroupTextMessage

		err = json.Unmarshal(data, &res)
		assert.Nil(t, err)

		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_REQUEST, res.Code)
		assert.EqualError(t, expectedError, res.Message)
		assert.Empty(t, res.Body)
	})

	mt.Run("HandleGroupConnectionsEP - Error author not exist", func(mt *mtest.T) {
		server.StartWebsocketService()

		expectedEmail := "jorge@mail.com"
		expectedGroupID := "6177226702-5T2de426p8arbt6sb4b128o63afaG9u3f-1727206726"
		expectedUsers := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID(), MockObjectID}

		messageBody := models.InboundGroupTextMessage{
			MessageID: "",
			AuthorID:  MockObjectID.Hex(),
			GroupID:   expectedGroupID,
			PushTokens: map[string]string{
				expectedUsers[0].Hex(): "333333333",
				expectedUsers[1].Hex(): "222222222",
				expectedUsers[2].Hex(): "111111111",
			},
			Body: "Hey everyone",
		}

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		returnedGroup := models.Group{
			ID:           primitive.NewObjectID(),
			GroupID:      expectedGroupID,
			Name:         "Some group name",
			Description:  "Test group",
			Participants: expectedUsers,
			Admins:       []primitive.ObjectID{MockObjectID},
			ProfileImage: "",
		}

		db := DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, false, nil
			},
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return &returnedGroup, nil
			},
			InsertGroupMessageDBMockFun: func(m any) (string, error) {
				return "", nil
			},
		}

		m := &media.MediaMock{}

		S := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleGroupConnectionsEP, &db, m)))
		defer S.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?gi=%s&ui=%s", strings.ReplaceAll(S.URL, "http", "ws"), expectedGroupID, expectedEmail), nil)
		assert.Nil(t, err)
		defer conn.Close()

		data2Send, err := json.Marshal(messageBody)
		assert.Nil(t, err)

		err = conn.WriteMessage(websocket.TextMessage, data2Send)
		assert.Nil(t, err)

		_, data, _ := conn.ReadMessage()

		var res models.OutboundGroupTextMessage

		err = json.Unmarshal(data, &res)
		assert.Nil(t, err)

		assert.True(t, res.Error)
		assert.Equal(t, server.NOT_ALLOWED, res.Code)
		assert.EqualError(t, errors.New("forbidden"), res.Message)
		assert.Empty(t, res.Body)

	})

	mt.Run("HandleGroupConnectionsEP - Error Not group found", func(mt *mtest.T) {

		server.StartWebsocketService()

		expectedError := mongo.ErrNoDocuments

		expectedEmail := "jorge@mail.com"
		expectedGroupID := ""
		expectedUsers := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID(), MockObjectID}

		messageBody := models.InboundGroupTextMessage{
			MessageID: "",
			AuthorID:  MockObjectID.Hex(),
			GroupID:   expectedGroupID,
			PushTokens: map[string]string{
				expectedUsers[0].Hex(): "333333333",
				expectedUsers[1].Hex(): "222222222",
				expectedUsers[2].Hex(): "111111111",
			},
			Body: "Hey everyone",
		}

		returnedUser := models.User{
			ID:    MockObjectID,
			Name:  "George",
			Email: expectedEmail,
			Credentials: models.UserCredentials{
				PushToken: "22222222",
			},
		}

		returnedGroup := models.Group{
			ID:           primitive.NewObjectID(),
			GroupID:      expectedGroupID,
			Name:         "Some group name",
			Description:  "Test group",
			Participants: expectedUsers,
			Admins:       []primitive.ObjectID{MockObjectID},
			ProfileImage: "",
		}

		db := DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			FindUserMockFunc: func(s string) (models.User, bool, error) {
				return returnedUser, true, nil
			},
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return &returnedGroup, expectedError
			},
			InsertGroupMessageDBMockFun: func(m any) (string, error) {
				return "", nil
			},
		}

		m := &media.MediaMock{}

		S := httptest.NewServer(http.HandlerFunc(decorators.HandlerWProvidersDecorator(HandleGroupConnectionsEP, &db, m)))
		defer S.Close()

		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws?gi=%s&ui=%s", strings.ReplaceAll(S.URL, "http", "ws"), expectedGroupID, expectedEmail), nil)
		assert.Nil(t, err)
		defer conn.Close()

		data2Send, err := json.Marshal(messageBody)
		assert.Nil(t, err)

		err = conn.WriteMessage(websocket.TextMessage, data2Send)
		assert.Nil(t, err)

		_, data, err := conn.ReadMessage()
		assert.Nil(t, err)

		var res models.OutboundGroupTextMessage

		err = json.Unmarshal(data, &res)
		assert.Nil(t, err)

		assert.True(t, res.Error)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.EqualError(t, expectedError, res.Message)
		assert.Empty(t, res.Body)

	})
}
