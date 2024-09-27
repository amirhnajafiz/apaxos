package database

import (
	"fmt"

	"github.com/f24-cse535/apaxos/internal/config/storage"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database is a module that handles database queries.
type Database struct {
	db *gorm.DB
}

// NewDatabase opens a MySQL connection and returns an instance of
// database struct.
func NewDatabase(cfg storage.MySQLConfig) (*Database, error) {
	// creating the MySQL dns from storage configs
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Database)

	// open connection to db
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %v", err)
	}

	return &Database{db: db}, nil
}
