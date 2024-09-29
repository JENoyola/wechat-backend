package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"wechat-back/internals/decorators"
	"wechat-back/internals/models"
	"wechat-back/internals/server"
	"wechat-back/providers/media"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestCreateNewGroupEP Tests the handler CreateNewGroupEP
func TestCreateNewGroupEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("CreateNewGroupEP - Success insertion", func(mt *mtest.T) {
		imageData := "iVBORw0KGgoAAAANSUhEUgAAAA0AAAANCAYAAABy6+R8AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAEOSURBVChTTZFRYsYgCIMR293/YefYAXa1KuwLavfTxlohQLD9/nynpWF8Mta21ftvdWgWyYbXu/E0wc0LzVghQVPwJpRxpLTusLvYPJWe/amizAdjTssRNsc0z94sOhFE9lzEVBs6oqrwFoOQT9AecQJ9WXO3i3ZPqWoVdJHxWUfGV4cEg6WIDtGvrXFrlTZVk+Z+c4bfUxkUXFlprfYECNIpk2ZVqxhhpglv468AbfePku19ER/uZkTUhCRFspdb6xJX/1unzG1yGBqr2YMj5NREL3+nWvErkwUFvCk6FpGi5UuiDhpD+jT5feQwIclQGTnWt7DFf1oNIlrlBlKj1Osy8xAUdN9cxbWTNPsDd7F5rPBhKbMAAAAASUVORK5CYII="

		data := models.Group{
			Name:         "Pirates fans",
			Description:  "Fan group of pirates. We share stories and more...",
			Participants: []primitive.ObjectID{MockObjectID},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		r, err := json.Marshal(data)
		assert.Nil(t, err)

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			InsertGroupDBMockFunc: func(g models.Group) (string, error) {
				return "", nil
			},
		}

		m := &media.MediaMock{}

		var body bytes.Buffer

		// add fields
		writter := multipart.NewWriter(&body)

		part, err := writter.CreateFormFile("avatar", "test.jpg")
		assert.Nil(t, err)

		d, err := base64.StdEncoding.DecodeString(imageData)
		assert.Nil(t, err)

		_, err = io.Copy(part, bytes.NewReader(d))
		assert.Nil(t, err)

		err = writter.WriteField("data", string(r))
		assert.Nil(t, err)

		err = writter.Close()
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/cg", &body)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", writter.FormDataContentType())

		rr := httptest.NewRecorder()

		handler := decorators.HandlerWProvidersDecorator(CreateNewGroupEP, db, m)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.COMPLETED, res.Code)
		assert.NotNil(t, res.DATA)

	})

	mt.Run("CreateNewGroupEP - Bad JSON", func(mt *mtest.T) {

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			InsertGroupDBMockFunc: func(g models.Group) (string, error) {
				return "", nil
			},
		}

		media := &media.MediaMock{}

		var body bytes.Buffer

		// add fields
		writter := multipart.NewWriter(&body)
		err := writter.WriteField("data", string("not a json"))
		assert.Nil(t, err)

		err = writter.Close()
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/cg", &body)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", writter.FormDataContentType())

		rr := httptest.NewRecorder()

		handler := decorators.HandlerWProvidersDecorator(CreateNewGroupEP, db, media)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_FIELD, res.Code)
		assert.Nil(t, res.DATA)

	})
	mt.Run("CreateNewGroupEP - Error inserting", func(mt *mtest.T) {
		imageData := "iVBORw0KGgoAAAANSUhEUgAAAA0AAAANCAYAAABy6+R8AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAEOSURBVChTTZFRYsYgCIMR293/YefYAXa1KuwLavfTxlohQLD9/nynpWF8Mta21ftvdWgWyYbXu/E0wc0LzVghQVPwJpRxpLTusLvYPJWe/amizAdjTssRNsc0z94sOhFE9lzEVBs6oqrwFoOQT9AecQJ9WXO3i3ZPqWoVdJHxWUfGV4cEg6WIDtGvrXFrlTZVk+Z+c4bfUxkUXFlprfYECNIpk2ZVqxhhpglv468AbfePku19ER/uZkTUhCRFspdb6xJX/1unzG1yGBqr2YMj5NREL3+nWvErkwUFvCk6FpGi5UuiDhpD+jT5feQwIclQGTnWt7DFf1oNIlrlBlKj1Osy8xAUdN9cxbWTNPsDd7F5rPBhKbMAAAAASUVORK5CYII="

		expectedError := errors.New("error trying to insert document")

		data := models.Group{
			Name:         "Pirates fans",
			Description:  "Fan group of pirates. We share stories and more...",
			Participants: []primitive.ObjectID{MockObjectID},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		r, err := json.Marshal(data)
		assert.Nil(t, err)

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			InsertGroupDBMockFunc: func(g models.Group) (string, error) {
				return "", expectedError
			},
		}

		media := &media.MediaMock{}

		var body bytes.Buffer

		// add fields
		writter := multipart.NewWriter(&body)

		part, err := writter.CreateFormFile("avatar", "test.jpg")
		assert.Nil(t, err)

		d, err := base64.StdEncoding.DecodeString(imageData)
		assert.Nil(t, err)

		_, err = io.Copy(part, bytes.NewReader(d))
		assert.Nil(t, err)

		err = writter.WriteField("data", string(r))
		assert.Nil(t, err)

		err = writter.Close()
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/cg", &body)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", writter.FormDataContentType())

		rr := httptest.NewRecorder()

		handler := decorators.HandlerWProvidersDecorator(CreateNewGroupEP, db, media)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)

	})
}

// TestUpdateGroupInfoEP tests the handler UpdateGroupInfoEP
func TestUpdateGroupInfoEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("UpdateGroupInfoEP - Success update", func(mt *mtest.T) {

		data := &models.Group{
			ID:          MockObjectID,
			GroupID:     "123456789",
			Name:        "Pirates fans",
			Description: "Fan group of pirates. We share stories and more...",
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return data, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		bod, err := json.Marshal(data)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/cg", bytes.NewReader(bod))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupInfoEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.OK, res.Code)
		assert.NotNil(t, res.DATA)

	})

	mt.Run("UpdateGroupInfoEP - bad json", func(mt *mtest.T) {

		data := "not a json"

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return nil, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		bod, err := json.Marshal(data)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/cg", bytes.NewReader(bod))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupInfoEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_FIELD, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupInfoEP - group not found", func(mt *mtest.T) {

		data := &models.Group{
			ID:          MockObjectID,
			GroupID:     "123456789",
			Name:        "Pirates fans",
			Description: "Fan group of pirates. We share stories and more...",
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return nil, mongo.ErrNoDocuments
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		bod, err := json.Marshal(data)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/cg", bytes.NewReader(bod))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupInfoEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_REQUEST, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupInfoEP - Error find one", func(mt *mtest.T) {

		expectedErr := errors.New("could not complete operation")

		data := &models.Group{
			ID:          MockObjectID,
			GroupID:     "123456789",
			Name:        "Pirates fans",
			Description: "Fan group of pirates. We share stories and more...",
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return data, expectedErr
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		bod, err := json.Marshal(data)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/cg", bytes.NewReader(bod))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupInfoEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupInfoEP - Error updating database", func(mt *mtest.T) {

		expectedError := errors.New("could not update document")

		data := &models.Group{
			ID:          MockObjectID,
			GroupID:     "123456789",
			Name:        "Pirates fans",
			Description: "Fan group of pirates. We share stories and more...",
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return data, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return expectedError
			},
		}

		bod, err := json.Marshal(data)
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "/cg", bytes.NewReader(bod))
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupInfoEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)

	})

}

// TestUpdateGroupAdminsEP test the handler UpdateGroupAdminsEP
func TestUpdateGroupAdminsEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("UpdateGroupAdminsEP - Success admin adition", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()

		operationType := OPERATION_ADD
		groupID := "123456789"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      groupID,
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&ads=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupAdminsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.OK, res.Code)
		assert.NotNil(t, res.DATA)
	})

	mt.Run("UpdateGroupAdminsEP - Success admin removal", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()

		operationType := OPERATION_REMOVE
		groupID := "123456789"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      groupID,
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID, secondUser},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&ads=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupAdminsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.OK, res.Code)
		assert.NotNil(t, res.DATA)
	})

	mt.Run("UpdateGroupAdminsEP - error bad group id", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()
		expectedError := mongo.ErrNoDocuments

		operationType := OPERATION_REMOVE
		groupID := "not a group id"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID, secondUser},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, expectedError
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&ads=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupAdminsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_REQUEST, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupAdminsEP - error findOne error", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()
		expectedError := errors.New("an error ocurred during find one operation")

		operationType := OPERATION_REMOVE
		groupID := "not a group id"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID, secondUser},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, expectedError
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&ads=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupAdminsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupAdminsEP - error no targets", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()
		expectedError := errors.New("no targets")

		operationType := OPERATION_REMOVE
		groupID := "123456789"
		targets := ""
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID, secondUser},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&ads=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupAdminsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
		assert.EqualError(t, expectedError, res.Message)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_FIELD, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupAdminsEP - error invalid primitive id", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()

		operationType := OPERATION_REMOVE
		groupID := "123456789"
		targets := "not a valid primitive id"
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID, secondUser},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&ads=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupAdminsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_CREDENTIALS, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupAdminsEP - error database update", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()
		expectedError := errors.New("could not update document")

		operationType := OPERATION_REMOVE
		groupID := "123456789"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID, secondUser},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return expectedError
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&ads=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupAdminsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedError, res.Message)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)

	})

}

// TestUpdateGroupParticipantsEP test the handler UpdateGroupParticipantsEP
func TestUpdateGroupParticipantsEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("UpdateGroupParticipantsEP - Success admin adition", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()

		operationType := OPERATION_ADD
		groupID := "123456789"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      groupID,
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&usrs=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupParticipantsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.OK, res.Code)
		assert.NotNil(t, res.DATA)
	})

	mt.Run("UpdateGroupParticipantsEP - Success admin removal", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()

		operationType := OPERATION_REMOVE
		groupID := "123456789"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      groupID,
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&usrs=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupParticipantsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.OK, res.Code)
		assert.NotNil(t, res.DATA)
	})

	mt.Run("UpdateGroupParticipantsEP - error bad group id", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()
		expectedError := mongo.ErrNoDocuments

		operationType := OPERATION_REMOVE
		groupID := "not a group id"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, expectedError
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&usrs=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupParticipantsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_REQUEST, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupParticipantsEP - error findOne error", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()
		expectedError := errors.New("an error ocurred during find one operation")

		operationType := OPERATION_REMOVE
		groupID := "not a group id"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, expectedError
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&usrs=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupParticipantsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupParticipantsEP - error no targets", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()
		expectedError := errors.New("no targets")

		operationType := OPERATION_REMOVE
		groupID := "123456789"
		targets := ""
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&usrs=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupParticipantsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
		assert.EqualError(t, expectedError, res.Message)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_FIELD, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupParticipantsEP - error invalid primitive id", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()

		operationType := OPERATION_REMOVE
		groupID := "123456789"
		targets := "not a valid primitive id"
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&usrs=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupParticipantsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_CREDENTIALS, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("UpdateGroupParticipantsEP - error database update", func(mt *mtest.T) {

		secondUser := primitive.NewObjectID()
		expectedError := errors.New("could not update document")

		operationType := OPERATION_REMOVE
		groupID := "123456789"
		targets := secondUser.Hex()
		adminName := "jorge"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      "123456789",
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID, secondUser},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			UpdateGroupDBMockFunc: func(m map[string]any, oi primitive.ObjectID) error {
				return expectedError
			},
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/uchta?ai=%s&ot=%s&gi=%s&usrs=%s", adminName, operationType, groupID, targets), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(UpdateGroupParticipantsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		// DO ASSERTION

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedError, res.Message)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)

	})

}

// TestDeleteGroupEP test handler DeleteGroupEP
func TestDeleteGroupEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("DeleteGroupEP - Success deletition", func(mt *mtest.T) {

		groupID := "123456"

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      groupID,
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			DeleteGroupDBMockFunc: func(s string) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/dgp?gi=%s", groupID), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(DeleteGroupEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusContinue, rr.Code)
		assert.False(t, res.Error)
		assert.Equal(t, server.COMPLETED, res.Code)
		assert.NotNil(t, res.DATA)

	})

	mt.Run("DeleteGroupEP - Error no target", func(mt *mtest.T) {

		groupID := ""

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      groupID,
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			DeleteGroupDBMockFunc: func(s string) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/dgp?gi=%s", groupID), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(DeleteGroupEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
		assert.True(t, res.Error)
		assert.Equal(t, server.BAD_FIELD, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("DeleteGroupEP - Error no group found", func(mt *mtest.T) {
		groupID := "123456"

		expectedErr := mongo.ErrNoDocuments

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return nil, expectedErr
			},
			DeleteGroupDBMockFunc: func(s string) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/dgp?gi=%s", groupID), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(DeleteGroupEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedErr, res.Message)
		assert.Equal(t, server.BAD_REQUEST, res.Code)
		assert.Nil(t, res.DATA)

	})

	mt.Run("DeleteGroupEP - Error deleting group", func(mt *mtest.T) {

		groupID := "123456"

		expectedError := errors.New("could not delete group")

		DBDoc := &models.Group{
			ID:           MockObjectID,
			GroupID:      groupID,
			Name:         "Wise Wizards",
			Description:  "Group about wizards",
			Participants: []primitive.ObjectID{MockObjectID},
			Admins:       []primitive.ObjectID{MockObjectID},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			GetGroupDBMockFunc: func(s string) (*models.Group, error) {
				return DBDoc, nil
			},
			DeleteGroupDBMockFunc: func(s string) error {
				return expectedError
			},
		}

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/dgp?gi=%s", groupID), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(DeleteGroupEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedError, res.Message)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)

	})

}

// TestSearchGroupsEP tests handler SearchGroupsEP
func TestSearchGroupsEP(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("SearchGroupsEP - Successful search", func(mt *mtest.T) {

		groupID := "1245678"

		results := []*models.Group{
			{
				ID:           MockObjectID,
				GroupID:      groupID,
				Name:         "Wise Wizards",
				Description:  "Group about wizards",
				Participants: []primitive.ObjectID{MockObjectID},
				Admins:       []primitive.ObjectID{MockObjectID},
			},
			{
				ID:           MockObjectID,
				GroupID:      groupID,
				Name:         "Pets mexico",
				Description:  "Group about wizards",
				Participants: []primitive.ObjectID{MockObjectID},
				Admins:       []primitive.ObjectID{MockObjectID},
			},
			{
				ID:           MockObjectID,
				GroupID:      groupID,
				Name:         "Cooking recipies",
				Description:  "Group about wizards",
				Participants: []primitive.ObjectID{MockObjectID},
				Admins:       []primitive.ObjectID{MockObjectID},
			},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			SearchGroupsMockFunc: func(i int, s string) ([]*models.Group, error) {

				return results, nil
			},
		}

		pg := 1
		query := ""

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/sgps?pg=%d&q=%s", pg, query), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(SearchGroupsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.False(t, res.Error)
		assert.NotNil(t, res.DATA)

	})

	mt.Run("SearchGroupsEP - Error no documents", func(mt *mtest.T) {
		expectedError := mongo.ErrNilCursor

		results := []*models.Group{}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			SearchGroupsMockFunc: func(i int, s string) ([]*models.Group, error) {

				return results, expectedError
			},
		}

		pg := 1
		query := ""

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/sgps?pg=%d&q=%s", pg, query), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(SearchGroupsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedError, res.Message)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)
	})

	mt.Run("SearchGroupsEP - Error database", func(mt *mtest.T) {
		groupID := "1245678"
		expectedError := errors.New("there was an error with the cursor")

		results := []*models.Group{
			{
				ID:           MockObjectID,
				GroupID:      groupID,
				Name:         "Wise Wizards",
				Description:  "Group about wizards",
				Participants: []primitive.ObjectID{MockObjectID},
				Admins:       []primitive.ObjectID{MockObjectID},
			},
			{
				ID:           MockObjectID,
				GroupID:      groupID,
				Name:         "Pets mexico",
				Description:  "Group about wizards",
				Participants: []primitive.ObjectID{MockObjectID},
				Admins:       []primitive.ObjectID{MockObjectID},
			},
			{
				ID:           MockObjectID,
				GroupID:      groupID,
				Name:         "Cooking recipies",
				Description:  "Group about wizards",
				Participants: []primitive.ObjectID{MockObjectID},
				Admins:       []primitive.ObjectID{MockObjectID},
			},
		}

		db := &DBMock{
			Client:       mt.Client,
			DatabaseName: MockDBName,
			SearchGroupsMockFunc: func(i int, s string) ([]*models.Group, error) {

				return results, expectedError
			},
		}

		pg := 1
		query := ""

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/sgps?pg=%d&q=%s", pg, query), nil)
		assert.Nil(t, err)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := decorators.HandlerDecorator(SearchGroupsEP, db)
		handler.ServeHTTP(rr, req)

		var res models.ServerResponse

		err = json.NewDecoder(rr.Body).Decode(&res)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.True(t, res.Error)
		assert.EqualError(t, expectedError, res.Message)
		assert.Equal(t, server.DB_ERROR, res.Code)
		assert.Nil(t, res.DATA)
	})

}
