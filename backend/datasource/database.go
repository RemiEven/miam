package datasource

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Necessary to load sqlite3 driver
	"github.com/sirupsen/logrus"
)

// databaseHolder holds a database connection
type databaseHolder struct {
	db *sql.DB
}

// newDatabaseHolder returns a new database holder
func newDatabaseHolder() (*databaseHolder, error) {
	db, err := sql.Open("sqlite3", "./miam.db")
	if err != nil {
		return nil, err
	}
	return &databaseHolder{db}, nil
}

// Close closes the connection to the database
func (holder *databaseHolder) Close() error {
	logrus.Info("Closing sqlite database connection")
	return holder.db.Close()
}
