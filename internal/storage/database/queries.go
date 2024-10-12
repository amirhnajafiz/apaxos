package database

import (
	"context"
	"fmt"

	"github.com/f24-cse535/apaxos/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
)

// InsertBlocks accepts an array of Block model and inserts them
// inside MongoDB.
func (d *Database) InsertBlocks(instances []*models.Block) error {
	ctx := context.TODO()

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
	ctx := context.TODO()

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
