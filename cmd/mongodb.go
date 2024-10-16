package cmd

import (
	"context"

	"github.com/f24-cse535/apaxos/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB command is used to check the connection to mongodb cluster.
type MongoDB struct {
	Cfg config.Config
}

func (m MongoDB) Main() error {
	// open a new connection to MongoDB cluster
	conn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.Cfg.MongoDB.URI))
	if err != nil {
		return err
	}

	// send a ping to confirm a successful connection
	if err := conn.Database(m.Cfg.MongoDB.Database).RunCommand(context.TODO(), bson.D{primitive.E{Key: "ping", Value: 1}}).Err(); err != nil {
		return err
	}

	return nil
}
