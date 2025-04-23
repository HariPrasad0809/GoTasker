package handlers

import (
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB initializes the database connection for the handlers package
func SetDB(database *gorm.DB) {
	db = database
}

// InitDB is an alias for SetDB to maintain compatibility with test files
func InitDB(database *gorm.DB) {
	SetDB(database)
}

// IsDBInitialized checks if the database connection is initialized
func IsDBInitialized() bool {
	return db != nil
}

// GetDB returns the current database instance (optional, use if needed)
func GetDB() *gorm.DB {
	return db
}
