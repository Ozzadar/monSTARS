package db

import (
	r "github.com/dancannon/gorethink"
	"github.com/ozzadar/monSTARS/models"
)

/*NewTransaction registers user in database.
- Username must be unique
- Email must be unique
*/
func NewTransaction(transaction *models.Transaction) (bool, string) {

	res, err := r.Table(TransactionsDB).Insert(transaction).RunWrite(Session)

	if err != nil {
		if res.Inserted == 0 {
			return false, "Transaction must be unique"
		}
		return false, "Error occurred"
	}

	return true, "User created successfully"
}
