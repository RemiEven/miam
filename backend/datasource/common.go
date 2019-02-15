package datasource

import (
	"strconv"

	"github.com/RemiEven/miam/common"
)

type sqliteID = int

// toSqliteID parses an id or return an error
func toSqliteID(ID string) (sqliteID, error) {
	intID, err := strconv.Atoi(ID)
	if err != nil {
		return 0, common.ErrInvalidID // TODO: wrap root error
	}
	return intID, nil
}

// fromSqliteID serializes an id
func fromSqliteID(ID sqliteID) string {
	return strconv.Itoa(ID)
}
