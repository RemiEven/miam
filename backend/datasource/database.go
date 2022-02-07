package datasource

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Necessary to load sqlite3 driver
	"github.com/rs/zerolog/log"
)

// DatabaseHolder holds a database connection
type DatabaseHolder struct {
	db *sql.DB
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
	log.Info().Msg("closing sqlite database connection")
	return holder.db.Close()
}
