package database

import (
	"context"
	"fmt"

	"github.com/f24-cse535/apaxos/internal/config/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// These const values are the collection names in the MongoDB cluster.
const (
	historyCollection   = "history"
	datastoreCollection = "ds"
)

// Database is a module that uses mongo-driver library to handle MongoDB queries.
type Database struct {
	conn      *mongo.Client
	history   *mongo.Collection
	datastore *mongo.Collection
}

// NewDatabase opens a MySQL connection and returns an instance of
// database struct.
func NewDatabase(cfg mongodb.Config, prefix string) (*Database, error) {
	// open a new connection to MongoDB cluster
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("[storage/database] failed to open a MongoDB connection: %v", err)
	}

	// create pointers to collections
	hpr := conn.Database(cfg.Database).Collection(fmt.Sprintf("%s_%s", prefix, historyCollection))
	dpr := conn.Database(cfg.Database).Collection(fmt.Sprintf("%s_%s", prefix, datastoreCollection))

	return &Database{
		conn:      conn,
		history:   hpr,
		datastore: dpr,
	}, nil
}
