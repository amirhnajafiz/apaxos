package database

import (
	"context"
	"fmt"

	"github.com/f24-cse535/apaxos/internal/config/storage"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is a module that uses mongo-driver library to handle MongoDB queries.
type Database struct {
	conn *mongo.Client
}

// NewDatabase opens a MySQL connection and returns an instance of
// database struct.
func NewDatabase(cfg storage.MongoDB) (*Database, error) {
	// open a new connection to MongoDB cluster
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("[storage/database] failed to open a MongoDB connection: %v", err)
	}

	return &Database{conn: conn}, nil
}
