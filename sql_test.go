package golangdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSQL(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// Execute a query
	query := "INSERT INTO customer(id, name) VALUES('admin','admin')"
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}

	fmt.Println("Executed query")
}

func TestQuerySQL(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// Execute a query
	query := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name string

		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("id: ", id)     // id: admin
		fmt.Println("name: ", name) // name: admin
	}
}

func TestQuerySQLComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// Execute a query
	query := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating sql.NullFloat64
		var birthDate, createdAt time.Time
		var married bool

		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("===================================================")
		fmt.Println("id: ", id)
		fmt.Println("name: ", name)
		fmt.Println("email: ", email.String)
		fmt.Println("balance: ", balance)
		fmt.Println("rating: ", rating.Float64)
		fmt.Println("birthDate: ", birthDate)
		fmt.Println("married: ", married)
		fmt.Println("createdAt: ", createdAt)
	}
}

func TestSQLInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "' OR TRUE; #"
	password := "admin"

	// Execute a query
	query := "SELECT username FROM user WHERE username = '" +
		username + "' AND password = '" + password + "' LIMIT 1"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		var username string

		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("username: ", username)
	} else {
		fmt.Println("User not found!")
	}
}

func TestQueryPreventSQLInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "' OR TRUE; #"
	password := "admin"

	// Execute a query
	query := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"

	rows, err := db.QueryContext(ctx, query, username, password)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		var username string

		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("username: ", username)
	} else {
		fmt.Println("User not found!")
	}
}

func TestExecPreventSQLInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin2"
	password := "admin2"

	// Execute a query
	query := "INSERT INTO user(username, password) VALUES(?, ?)"
	_, err := db.ExecContext(ctx, query, username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Executed query")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "admin@admin.com"
	comment := "test koment"

	// Execute a query
	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, query, email, comment)
	if err != nil {
		panic(err)
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println(insertedId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// Execute a query
	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for i := 0; i < 10; i++ {
		email := "admin" + strconv.Itoa(i) + "@admin.com"
		comment := "test koment nomor " + strconv.Itoa(i)

		result, err := stmt.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("comment id: ", id)
	}
}

func TestDatabaseTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// Execute a query
	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		email := "admin" + strconv.Itoa(i) + "@admin.com"
		comment := "test koment nomor " + strconv.Itoa(i)

		result, err := stmt.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("comment id: ", id)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	// err = tx.Rollback()
	// if err != nil {
	// 	panic(err)
	// }
}
