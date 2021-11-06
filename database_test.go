package golangdatabase

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestOpenDatabaseConnection(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/belajar_golang_database")
	if err != nil {
		t.Errorf("Error opening database connection: %s", err)
		panic(err)
	}

	defer db.Close()
}
