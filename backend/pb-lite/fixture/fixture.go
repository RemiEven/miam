package fixture

import (
	"fmt"

	"github.com/remieven/miam/datasource"
)

// PrepareDatabase can be used to prepare a database by executing a list of SQL statements
func PrepareDatabase(statements ...string) func(*datasource.DatabaseHolder) error {
	return func(holder *datasource.DatabaseHolder) error {
		for index, statement := range statements {
			if _, err := holder.DB.Exec(statement); err != nil {
				return fmt.Errorf("failed to execute statement %v: %w", index, err)
			}
		}
		return nil
	}
}
