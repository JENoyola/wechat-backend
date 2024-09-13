package database

import "go.mongodb.org/mongo-driver/bson/primitive"

// Mock environment variables
const (
	mockURL         = "mongodb://localhost:27017"
	mockDBName      = "test_db"
	ObjectIDMockHex = "66d6561e43416dd7f7eb6aa4"
)

var (
	ObjectIDMock = primitive.NewObjectID()
)

// MockDB holds a mock of the DB struct for testing
type MockDB struct {
	Database string
}

// Mock the StartDatabase function to return a mocked DB object
func StartMockDatabase() *DB {
	// Instead of connecting to a real DB, return a mock DB
	return &DB{
		client:   nil, // No real client
		Database: mockDBName,
	}
}
