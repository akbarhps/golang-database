# Golang Database

Sumber Tutorial:
[Udemy](https://www.udemy.com/course/pemrograman-go-lang-pemula-sampai-mahir/learn/lecture/25028914#questions) |
[Slide](https://docs.google.com/presentation/d/15pvN3L3HTgA9aIMNkm03PzzIwlff0WDE6hOWWut9pg8/edit#slide=id.p)


## Pengenalan Package Database
---

- Bahasa pemrograman Go-Lang secara default memiliki sebuah package bernama database
- Package database adalah package yang berisikan kumpulan standard interface yang menjadi standard untuk berkomunikasi ke database
- Hal ini menjadikan kode program yang kita buat untuk mengakses jenis database apapun bisa menggunakan kode yang sama
- Yang berbeda hanya kode SQL yang perlu kita gunakan sesuai dengan database yang kita gunakan


### Cara Kerja Package Database

![Cara Kerja Package Database](https://user-images.githubusercontent.com/69947442/140595539-28452b16-e23a-4eb5-b6d9-cb0a69e59e46.png)


### MySQL

- Pada materi kali ini kita akan menggunakan MySQL sebagai Database Management System
- Jadi pastikan teman-teman sudah mengerti tentang MySQL


### Menambahkan Database Driver
---

### Database Driver

- Sebelum kita membuat kode program menggunakan database di Go-Lang, terlebih dahulu kita wajib menambahkan driver database nya
- Tanpa driver database, maka package database di Go-Lang tidak mengerti apapun, karena hanya berisi kontrak interface saja
- https://golang.org/s/sqldrivers 


### Menambahkan Module Database MySQL

```bash
go get -u github.com/go-sql-driver/mysql
```


### Import Package MySQL

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)
```


## Membuat Koneksi Database
---

- Hal yang pertama akan kita lakukan ketika membuat aplikasi yang akan menggunakan database adalah melakukan koneksi ke database nya
- Untuk melakukan koneksi ke databsae di Golang, kita bisa membuat object `sql.DB` menggunakan function `sql.Open(driver, dataSourceName)`
- Untuk menggunakan database MySQL, kita bisa menggunakan driver `"mysql"`
- Sedangkan untuk dataSourceName, tiap database biasanya punya cara penulisan masing-masing, misal di MySQL, kita bisa menggunakan dataSourceName seperti dibawah ini :
  - `username:password@tcp(host:port)/database_name`
- Jika object `sql.DB` sudah tidak digunakan lagi, disarankan untuk menutupnya menggunakan function `Close()`


### Kode: Membuat Koneksi Database

```go
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
```


## Database Pooling
---

- `sql.DB` di Golang sebenarnya bukanlah sebuah koneksi ke database
- Melainkan sebuah pool ke database, atau dikenal dengan konsep Database Pooling
- Di dalam `sql.DB`, Golang melakukan management koneksi ke database secara otomatis. Hal ini menjadikan kita tidak perlu melakukan management koneksi database secara manual
- Dengan kemampuan database pooling ini, kita bisa menentukan jumlah minimal dan maksimal koneksi yang dibuat oleh Golang, sehingga tidak membanjiri koneksi ke database, karena biasanya ada batas maksimal koneksi yang bisa ditangani oleh database yang kita gunakan


### Pengaturan Database Pooling

| Method                              | Keterangan                                                             |
| ----------------------------------- | ---------------------------------------------------------------------- |
| `(DB) SetMaxIdleConns(number)`      | Pengaturan berapa jumlah koneksi minimal yang dibuat                   |
| `(DB) SetMaxOpenConns(number)`      | Pengaturan berapa jumlah koneksi maksimal yang dibuat                  |
| `(DB) SetConnMaxIdleTime(duration)` | Pengaturan berapa lama koneksi yang sudah tidak digunakan akan dihapus |
| `(DB) SetConnMaxLifetime(duration)` | Pengaturan berapa lama koneksi boleh digunakan                         |


### Kode: Database Pooling

```go
func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/belajar_golang_database")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
```


## Eksekusi Perintah SQL
---

- Saat membuat aplikasi menggunakan database, sudah pasti kita ingin berkomunikasi dengan database menggunakan perintah SQL
- Di Golang juga menyediakan function yang bisa kita gunakan untuk mengirim perintah SQL ke database menggunakan function `(DB) ExecContext(context, sql, params)`
- Ketika mengirim perintah SQL, kita butuh mengirimkan context, dan seperti yang sudah pernah kita pelajari di course Golang Context, dengan context, kita bisa mengirim sinyal cancel jika kita ingin membatalkan pengiriman perintah SQL nya


### Kode: Membuat Table Customer

```sql
CREATE TABLE customer
(
    id   VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB;
```


### Kode: Mengirim Perintah SQL Insert

```go
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
```


## Query SQL
---

- Untuk operasi SQL yang tidak membutuhkan hasil, kita bisa menggunakan perintah Exec, namun jika kita membutuhkan result, seperti SELECT SQL, kita bisa menggunakan function yang berbeda
- Function untuk melakukan query ke database, bisa menggunakan function `(DB) QueryContext(context, sql, params)`


### Kode: Query SQL

```go
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
}
```


### Rows

- Hasil Query function adalah sebuah data structs `sql.Rows`
- Rows digunakan untuk melakukan iterasi terhadap hasil dari query
- Kita bisa menggunakan function `(Rows) Next() (boolean)` untuk melakukan iterasi terhadap data hasil query, jika return data false, artinya sudah tidak ada data lagi didalam result
- Untuk membaca tiap data, kita bisa menggunakan `(Rows) Scan(columns...)`
- Dan jangan lupa, setelah menggunakan Rows, jangan lupa untuk menutupnya menggunakan `(Rows) Close()`


### Kode: Rows

```go
for rows.Next() {
  var id, name string

  err := rows.Scan(&id, &name)
  if err != nil {
    panic(err)
  }

  fmt.Println("id: ", id)     // id: admin
  fmt.Println("name: ", name) // name: admin
}
```


## Tipe Data Column
---

- Sebelumnya kita hanya membuat table dengan tipe data di kolom nya berupa VARCHAR
- Untuk VARCHAR di database, biasanya kita gunakan String di Golang
- Bagaimana dengan tipe data yang lain?
- Apa representasinya di Golang, misal tipe data timestamp, date dan lain-lain
 

### Kode: Alter Table Customer

```sql 
ALTER TABLE customer
    ADD COLUMN email      varchar(100),
    ADD COLUMN balance    INT       DEFAULT 0,
    ADD COLUMN rating     DOUBLE    DEFAULT 0.0,
    ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN birth_date DATE,
    ADD COLUMN married    BOOLEAN   DEFAULT FALSE;
```


### Mapping Tipe Data

| Database                        | Golang           |
| ------------------------------- | ---------------- |
| VARCHAR, CHAR                   | string           |
| INT, BIGINT                     | int32, int64     |
| FLOAT, DOUBLE                   | float32, float64 |
| BOOLEAN                         | bool             |
| DATE, DATETIME, TIME, TIMESTAMP | time.Time        |


### Kode: Insert Data Customer

```sql
INSERT INTO customer(id, name, email, balance, rating, birth_date, married)
VALUES ('admin', 'admin', 'admin@admin.com', 100000, 5.0, '1999-9-9', true),
       ('admin2', 'admin2', 'admin2@admin.com', 100000, 5.0, '1999-9-9', true);
```


### Error Tipe Data Date

```bash
➜ go test -v -timeout 10s -run=TestQuerySQLComplex
=== RUN   TestQuerySQLComplex
--- FAIL: TestQuerySQLComplex (0.01s)
panic: sql: Scan error on column index 5, name "birth_date": unsupported Scan, storing driver.Value type []uint8 into type *time.Time [recovered]
        panic: sql: Scan error on column index 5, name "birth_date": unsupported Scan, storing driver.Value type []uint8 into type *time.Time

goroutine 19 [running]:
testing.tRunner.func1.2({0x10862a0, 0xc000096740})
        C:/Program Files/Go/src/testing/testing.go:1209 +0x24e
testing.tRunner.func1()
        C:/Program Files/Go/src/testing/testing.go:1212 +0x218
panic({0x10862a0, 0xc000096740})
        C:/Program Files/Go/src/runtime/panic.go:1038 +0x215
golang-database.TestQuerySQLComplex(0x0)
        C:/Users/Akbar/GolandProjects/src/golang-database/sql_test.go:78 +0x55d
testing.tRunner(0xc0000851e0, 0x10ca580)
        C:/Program Files/Go/src/testing/testing.go:1259 +0x102
created by testing.(*T).Run
        C:/Program Files/Go/src/testing/testing.go:1306 +0x35a
exit status 2
FAIL    golang-database 0.062s
```

- Secara default, Driver MySQL untuk Golang akan melakukan query tipe data DATE, DATETIME, TIMESTAMP menjadi []byte / []uint8. Dimana ini bisa dikonversi menjadi String, lalu di parsing menjadi time.Time
- Namun hal ini merepotkan jika dilakukan manual, kita bisa meminta Driver MySQL untuk Golang secara otomatis melakukan parsing dengan menambahkan parameter `parseTime=true`


### Kode: Get Connection

```go
func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/belajar_golang_database?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
```


### Kode: Query Ke Database Complex

```go
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
		var id, name, email string
		var balance int32
		var rating float32
		var birthDate, createdAt time.Time
		var married bool

		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Printf("id: %s, name: %s, email: %s, balance: %d, rating: %f, birth_date: %s, married: %t, created_at: %s\n\n", id, name, email, balance, rating, birthDate, married, createdAt)
	}
}
```


### Output

```bash
➜ go test -v -timeout 10s -run=TestQuerySQLComplex
=== RUN   TestQuerySQLComplex
id: admin, name: admin, email: admin@admin.com, balance: 100000, rating: 5.000000, birth_date: 1999-09-09 00:00:00 +0000 UTC, married: true, created_at: 2021-11-06 10:59:25 +0000 UTC

id: admin2, name: admin2, email: admin2@admin.com, balance: 100000, rating: 5.000000, birth_date: 1999-09-09 00:00:00 +0000 UTC, married: true, created_at: 2021-11-06 10:59:40 +0000 UTC

--- PASS: TestQuerySQLComplex (0.01s)
PASS
ok      golang-database 0.040s
```


### Nullable Type

- Golang database tidak mengerti dengan tipe data NULL di database
- Oleh karena itu, khusus untuk kolom yang bisa NULL di database, akan jadi masalah jika kita melakukan Scan secara bulat-bulat menggunakan tipe data representasinya di Golang


### Kode: Insert Data Null

```sql
INSERT INTO customer(id, name, email, balance, rating, birth_date, married)
VALUES ('admin', 'admin', 'null', 100000, null, '1999-9-9', true);

UPDATE customer
set email  = null,
    rating = null
WHERE id = 'admin';
```


### Error Data Null

```bash
➜ go test -v -timeout 10s -run=TestQuerySQLComplex
=== RUN   TestQuerySQLComplex
--- FAIL: TestQuerySQLComplex (0.01s)
panic: sql: Scan error on column index 2, name "email": converting NULL to string is unsupported [recovered]
        panic: sql: Scan error on column index 2, name "email": converting NULL to string is unsupported

goroutine 19 [running]:
testing.tRunner.func1.2({0x8e62a0, 0xc00018a160})
        C:/Program Files/Go/src/testing/testing.go:1209 +0x24e
testing.tRunner.func1()
        C:/Program Files/Go/src/testing/testing.go:1212 +0x218
panic({0x8e62a0, 0xc00018a160})
        C:/Program Files/Go/src/runtime/panic.go:1038 +0x215
golang-database.TestQuerySQLComplex(0x0)
        C:/Users/Akbar/GolandProjects/src/golang-database/sql_test.go:78 +0x55d
testing.tRunner(0xc0000851e0, 0x92a5c0)
        C:/Program Files/Go/src/testing/testing.go:1259 +0x102
created by testing.(*T).Run
        C:/Program Files/Go/src/testing/testing.go:1306 +0x35a
exit status 2
FAIL    golang-database 0.043s
```

- Konversi secara otomatis NULL tidak didukung oleh Driver MySQL Golang
- Oleh karena itu, khusus tipe kolom yang bisa NULL, kita perlu menggunakan tipe data yang ada dalam package sql


### Tipe Data Nullable

| Golang    | Nullable                   |
| --------- | -------------------------- |
| string    | `database/sql.NullString`  |
| bool      | `database/sql.NullBool`    |
| float64   | `database/sql.NullFloat64` |
| int32     | `database/sql.NullInt32`   |
| int64     | `database/sql.NullInt64`   |
| time.Time | `database/sql.NullTime`    |


### Kode: Nullable

```go
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
```


### Output

```bash
➜ go test -v -timeout 10s -run=TestQuerySQLComplex
=== RUN   TestQuerySQLComplex
===================================================
id:  admin
name:  admin
email:
balance:  100000
rating:  0
birthDate:  1999-09-09 00:00:00 +0000 UTC
married:  true
createdAt:  2021-11-06 10:59:25 +0000 UTC
===================================================
id:  admin2
name:  admin2
email:  null
balance:  100000
rating:  0
birthDate:  1999-09-09 00:00:00 +0000 UTC
married:  true
createdAt:  2021-11-06 10:59:40 +0000 UTC
--- PASS: TestQuerySQLComplex (0.01s)
PASS
ok      golang-database 0.039s
```


### Kode: Mengecek Null atau Tidak

```go
  if email.Valid {
    fmt.Println("email: ", email.String)
  }
  if rating.Valid {
    fmt.Println("rating: ", rating.Float64)
  }
```


### Output

```bash
➜ go test -v -timeout 10s -run=TestQuerySQLComplex
=== RUN   TestQuerySQLComplex
===================================================
id:  admin
name:  admin
balance:  100000
birthDate:  1999-09-09 00:00:00 +0000 UTC
married:  true
createdAt:  2021-11-06 10:59:25 +0000 UTC
===================================================
id:  admin2
name:  admin2
email:  null
balance:  100000
birthDate:  1999-09-09 00:00:00 +0000 UTC
married:  true
createdAt:  2021-11-06 10:59:40 +0000 UTC
--- PASS: TestQuerySQLComplex (0.01s)
PASS
ok      golang-database 0.043s
```


## SQL Injection
---

### SQL Dengan Parameter

- Saat membuat aplikasi, kita tidak mungkin akan melakukan hardcode perintah SQL di kode Golang kita
- Biasanya kita akan menerima input data dari user, lalu membuat perintah SQL dari input user, dan mengirimnya menggunakan perintah SQL


### Kode: Membuat Table User

```sql
CREATE TABLE user
(
    username VARCHAR(100) NOT NULL ,
    password VARCHAR(100) NOT NULL ,
    PRIMARY KEY (username)
) ENGINE = InnoDB;

INSERT INTO user(username, password)
VALUES('admin', 'admin');
```


### Kode: SQL Query Dengan Parameter

```go
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
```


### SQL Injection

- SQL Injection adalah sebuah teknik yang menyalahgunakan sebuah celah keamanan yang terjadi dalam lapisan basis data sebuah aplikasi.
- Biasa, SQL Injection dilakukan dengan mengirim input dari user dengan perintah yang salah, sehingga menyebabkan hasil SQL yang kita buat menjadi tidak valid
- SQL Injection sangat berbahaya, jika sampai kita salah membuat SQL, bisa jadi data kita tidak aman


### Kode: SQL Injection

```go
username := "' OR TRUE; #"
password := "any"
```


### Output

```bash
➜ go test -v -timeout 10s -run=TestSQLInject
=== RUN   TestSQLInjection
username:  admin
--- PASS: TestSQLInjection (0.01s)
PASS
ok      golang-database 0.040s
```


### Solusinya?

- Jangan membuat query SQL secara manual dengan menggabungkan String secara bulat-bulat
- Jika kita membutuhkan parameter ketika membuat SQL, kita bisa menggunakan function Execute atau Query dengan parameter yang akan kita bahas di chapter selanjutnya


## SQL Dengan Parameter
---

- Jangan membuat query SQL secara manual dengan menggabungkan String secara bulat-bulat
- Jika kita membutuhkan parameter ketika membuat SQL, kita bisa menggunakan function Execute atau Query dengan parameter yang akan kita bahas di chapter selanjutnya


### Contoh SQL

- SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1
- INSERT INTO user(username, password) VALUES (?, ?)
- Dan lain-lain


### Kode: Query Dengan Parameter

```go
username := "' OR TRUE; #"
password := "admin"

// Execute a query
query := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"

rows, err := db.QueryContext(ctx, query, username, password)
if err != nil {
  panic(err)
}
```


### Output

```bash
➜ go test -v -timeout 10s -run=TestQueryPreventSQLInjection
=== RUN   TestQueryPreventSQLInjection
User not found!
--- PASS: TestQueryPreventSQLInjection (0.01s)
PASS
ok      golang-database 0.049s
```


### Kode: Exec Dengan Parameter

```go
username := "admin2"
password := "admin2"

// Execute a query
query := "INSERT INTO user(username, password) VALUES(?, ?)"
_, err := db.ExecContext(ctx, query, username, password)
if err != nil {
  panic(err)
}
```


### Output

```bash
➜ go test -v -timeout 10s -run=TestExecPreventSQLInjection
=== RUN   TestExecPreventSQLInjection
Executed query
--- PASS: TestExecPreventSQLInjection (0.01s)
PASS
ok      golang-database 0.047s
```


## Auto Increment
---

- Kadang kita membuat sebuah table dengan id auto increment
- Dan kadang pula, kita ingin mengambil data id yang sudah kita insert ke dalam MySQL
- Sebenarnya kita bisa melakukan query ulang ke database menggunakan SELECT LAST_INSERT_ID()
- Tapi untungnya di Golang ada cara yang lebih mudah
- Kita bisa menggunakan function `(Result) LastInsertId()` untuk mendapatkan Id terakhir yang dibuat secara auto increment 
- Result adalah object yang dikembalikan ketika kita menggunakan function Exec


### Kode: Membuat Table

```sql
CREATE TABLE comments
(
    id INT NOT NULL  AUTO_INCREMENT,
    email VARCHAR(100) NOT NULL ,
    comment TEXT,
    PRIMARY KEY (id)
) ENGINE = InnoDB;
```


### Kode: (Result) LastInsertId()

```go
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
```


### Output

```bash
➜ go test -v -timeout 10s -run=TestAutoIncrement
=== RUN   TestAutoIncrement
1
--- PASS: TestAutoIncrement (0.02s)
PASS
ok      golang-database 0.047s
~\..\..\golang-database
➜ go test -v -timeout 10s -run=TestAutoIncrement
=== RUN   TestAutoIncrement
2
--- PASS: TestAutoIncrement (0.01s)
PASS
ok      golang-database 0.037s
~\..\..\golang-database
➜ go test -v -timeout 10s -run=TestAutoIncrement
=== RUN   TestAutoIncrement
3
--- PASS: TestAutoIncrement (0.01s)
PASS
ok      golang-database 0.034s
```


## Prepare Statement
---


### Query atau Exec Dengan Parameter

- Saat kita menggunakan Function Query atau Exec yang menggunakan parameter, sebenarnya implementasi dibawah nya menggunakan Prepare Statement
- Jadi tahapan pertama statement nya disiapkan terlebih dahulu, setelah itu baru di isi dengan parameter
- Kadang ada kasus kita ingin melakukan beberapa hal yang sama sekaligus, hanya berbeda parameternya. Misal insert data langsung banyak
- Pembuatan Prepare Statement bisa dilakukan dengan manual, tanpa harus mennggunakan Query atau Exec dengan parameter


### Prepare Statement

- Saat kita membuat Prepare Statement, secara otomatis akan mengenali koneksi database yang digunakan
- Sehingga ketika kita mengeksekusi Prepare Statement berkali-kali, maka akan menggunakan koneksi yang sama dan lebih efisien karena pembuatan prepare statement nya hanya sekali diawal saja
- Jika menggunakan Query dan Exec dengan parameter, kita tidak bisa menjamin bahwa koneksi yang digunakan akan sama, oleh karena itu, bisa jadi prepare statement akan selalu dibuat berkali-kali walaupun kita menggunakan SQL yang sama
- Untuk membuat Prepare Statement, kita bisa menggunakan function `(DB) Prepare(context, sql)`
- Prepare Statement direpresentasikan dalam struct `database/sql.Stmt`
- Sama seperti resource sql lainnya, Stmt harus di `Close()` jika sudah tidak digunakan lagi


### Kode: Membuat Statement

```go
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
```


### Output

```bash
➜ go test -v -timeout 10s -run=TestPrepareStatement
=== RUN   TestPrepareStatement
comment id:  4
comment id:  5
comment id:  6
comment id:  7
comment id:  8
comment id:  9
comment id:  10
comment id:  11
comment id:  12
comment id:  13
--- PASS: TestPrepareStatement (0.03s)
PASS
ok      golang-database 0.058s
```


## Database Transaction
---

- Salah satu fitur andalan di database adalah transaction
- Materi database transaction sudah saya bahas dengan tuntas di materi MySQL database, jadi silahkan pelajari di course tersebut
- Di course ini kita akan fokus bagaimana menggunakan database transaction di Golang


### Transaction di Golang

- Secara default, semua perintah SQL yang kita kirim menggunakan Golang akan otomatis di commit, atau istilahnya auto commit
- Namun kita bisa menggunakan fitur transaksi sehingga SQL yang kita kirim tidak secara otomatis di commit ke database
- Untuk memulai transaksi, kita bisa menggunakan function `(DB) Begin()`, dimana akan menghasilkan struct Tx yang merupakan representasi Transaction
- Struct Tx ini yang kita gunakan sebagai pengganti DB untuk melakukan transaksi, dimana hampir semua function di DB ada di Tx, seperti Exec, Query atau Prepare
- Setelah selesai proses transaksi, kita bisa gunakan function `(Tx) Commit()` untuk melakukan commit atau `Rollback()`


### Kode: Transaction di Golang

```go
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
```


## Repository Pattern
---

- Dalam buku Domain-Driven Design, Eric Evans menjelaskan bahwa “repository is a mechanism for encapsulating storage, retrieval, and search behavior, which emulates a collection of objects”
- Pattern Repository ini biasanya digunakan sebagai jembatan antar business logic aplikasi kita dengan semua perintah SQL ke database
- Jadi semua perintah SQL akan ditulis di Repository, sedangkan business logic kode program kita hanya cukup menggunakan Repository tersebut


### Diagram Repository Pattern

![Diagram Repository Pattern](https://user-images.githubusercontent.com/69947442/140599899-fa25ee32-3878-4191-8825-606c0aec64ae.png)


### Entity/Model

- Dalam pemrograman berorientasi object, biasanya sebuah tabel di database akan selalu dibuat representasinya sebagai class Entity atau Model, namun di Golang, karena tidak mengenal Class, jadi kita akan representasikan data dalam bentuk Struct
- Ini bisa mempermudah ketika membuat kode program
- Misal ketika kita query ke Repository, dibanding mengembalikan array, alangkah baiknya Repository melakukan konversi terlebih dahulu ke struct Entity / Model, sehingga kita tinggal menggunakan objectnya saja


### Kode: Struct Model/Entity

```go
type Comment struct {
	Id      int32
	Email   string
	Comment string
}
```


### Kode: Interface Repository

```go
type CommentRepository interface {
	Insert(ctx context.Context, comment *entity.Comment) (*entity.Comment, error)
	FindById(ctx context.Context, id int32) (*entity.Comment, error)
	FindAll(ctx context.Context) ([]*entity.Comment, error)
}
```


### Kode: Implementasi Repository

```go
type commentRepositoryImpl struct {
	DB *sql.DB
}

func (repo *commentRepositoryImpl) Insert(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := repo.DB.ExecContext(ctx, query, comment.Email, comment.Comment)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	comment.Id = int32(id)
	return comment, nil
}

func (repo *commentRepositoryImpl) FindById(ctx context.Context, id int32) (*entity.Comment, error) {
	query := "SELECT id, email, comment FROM comments WHERE id = ? LIMIT 1"
	rows, err := repo.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comment := entity.Comment{}
	if rows.Next() {
		err := rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		if err != nil {
			return nil, err
		}
	}

	return &comment, nil
}

func (repo *commentRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Comment, error) {
	query := "SELECT id, email, comment FROM comments"
	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []*entity.Comment{}
	for rows.Next() {
		comment := entity.Comment{}
		err := rows.Scan(&comment.Id, &comment.Email, &comment.Comment)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}
```


### Kode: Implementasi New Repository

```go
func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{DB: db}
}
```


### Kode: Test Comment Repository

```go
import (
	"context"
	"golang-database/entity"
	"testing"

	golang_database "golang-database"
)

func TestCommentInsert(t *testing.T) {
	ctx := context.Background()
	commentRepository := NewCommentRepository(golang_database.GetConnection())
	comment := &entity.Comment{
		Email:   "test@test.com",
		Comment: "test comment dengan repository",
	}

	comment, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}

	t.Logf("%+v", comment)
}

func TestCommentFindById(t *testing.T) {
	ctx := context.Background()
	commentRepository := NewCommentRepository(golang_database.GetConnection())
	comment, err := commentRepository.FindById(ctx, 1)
	if err != nil {
		panic(err)
	}

	t.Logf("%+v", comment)
}

func TestCommentFindByIdNotFound(t *testing.T) {
	ctx := context.Background()
	commentRepository := NewCommentRepository(golang_database.GetConnection())
	comment, err := commentRepository.FindById(ctx, 0)
	if err != nil {
		panic(err)
	}

	t.Logf("%+v", comment)
}

func TestCommentFindAll(t *testing.T) {
	ctx := context.Background()
	commentRepository := NewCommentRepository(golang_database.GetConnection())
	comments, err := commentRepository.FindAll(ctx)
	if err != nil {
		panic(err)
	}

	for _, comment := range comments {
		t.Logf("%+v", comment)
	}
}
```