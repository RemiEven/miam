package datasource

import (
	"fmt"
	"strconv"
)

type sqliteID = int

// toSqliteID parses an id or return an error
func toSqliteID(ID string) (sqliteID, error) {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return 0, fmt.Errorf("failed to parse as an int: %w", err)
	}
	return intID, nil
}

// fromSqliteID serializes an id
func fromSqliteID(ID sqliteID) string {
	return strconv.Itoa(ID)
}
