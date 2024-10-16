package database

import (
	"context"
	"fmt"

	"github.com/f24-cse535/apaxos/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertBlocks accepts an array of Block model and inserts them
// inside MongoDB.
func (d *Database) InsertBlocks(instances []*models.Block) error {
	ctx := context.Background()

	// convert models.Block to interface type
	var interfaceObjects []interface{}
	for _, obj := range instances {
		interfaceObjects = append(interfaceObjects, obj)
	}

	// use insert many to store blocks
	if _, err := d.history.InsertMany(ctx, interfaceObjects); err != nil {
		return fmt.Errorf("[storage/database] failed to insert objects: %v", err)
	}

	return nil
}

// GetBlocks returns a list the current committed blocks.
func (d *Database) GetBlocks() ([]*models.Block, error) {
	ctx := context.Background()

	// fetch all blocks
	cursor, err := d.history.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("[storage/database] failed to fetch all blocks: %v", err)
	}
	defer cursor.Close(ctx)

	// convert MongoDB objects into models.Block
	var blocks []*models.Block

	for cursor.Next(ctx) {
		var block models.Block
		if err := cursor.Decode(&block); err != nil {
			return nil, fmt.Errorf("[storage/database] failed to decode object: %v", err)
		}

		blocks = append(blocks, &block)
	}

	return blocks, nil
}

// InsertState is used to store a system snapshot into database.
func (d *Database) InsertState(instance *models.State) error {
	ctx := context.Background()

	if _, err := d.states.InsertOne(ctx, instance); err != nil {
		return fmt.Errorf("[storage/database] failed to store state: %v", err)
	}

	return nil
}

// GetLastState is used to retrieve the previous node state.
func (d *Database) GetLastState() (*models.State, error) {
	ctx := context.Background()

	// find the last inserted item by sorting `_id` in descending order and limiting to 1
	filter := bson.D{primitive.E{Key: "_id", Value: "-1"}}
	opts := options.FindOne().SetSort(filter)

	var state models.State
	err := d.states.FindOne(ctx, bson.D{}, opts).Decode(&state)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("[storage/database] failed to load snapshot: %v", err)
	}

	return &state, nil
}
