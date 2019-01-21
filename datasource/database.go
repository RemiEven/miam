package datasource

import (
	"database/sql"
	"log"

	// Necessary to load sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// DatabaseHolder holds a database connection
type DatabaseHolder struct {
	db *sql.DB
}

// NewDatabaseHolder returns a new database holder
func NewDatabaseHolder() (*DatabaseHolder, error) {
	db, err := sql.Open("sqlite3", "./miam.db")
	if err != nil {
		return nil, err
	}
	return &DatabaseHolder{db}, nil
}

// Close closes the connection to the database
func (holder *DatabaseHolder) Close() error {
	log.Println("Closing database connection")
	return holder.db.Close()
}
