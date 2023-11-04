package datasource

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3" // Necessary to load sqlite3 driver
)

// DatabaseHolder holds a database connection
type DatabaseHolder struct {
	DB *sql.DB
}

// NewDatabaseHolder returns a new database holder
func NewDatabaseHolder(dbFilePath string) (*DatabaseHolder, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}
	return &DatabaseHolder{db}, nil
}

// Close cleanly closes the connection to the database
func (holder *DatabaseHolder) Close() error {
	slog.Info("closing sqlite database connection")
	return holder.DB.Close()
}
