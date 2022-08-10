package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	usingHex             = flag.Bool("hex", false, "use hex encoding for the binary result of query")
	usingMysql           = flag.Bool("mysql", false, "use mysql database")
	createTable          = flag.Bool("create", false, "create table in the database")
	dropTable            = flag.Bool("drop", false, "Drop current table from the database")
	insertRandomValues   = flag.Int("insert", 0, "insert N random values into the database")
	dbport               = flag.Int("dbport", 9393, "db port")
	poisonRecordToInsert = flag.String("insert_poison", "", "insert poison record (should be in BASE64 format)")
	selectAllFromTable   = flag.Int("select", 0, "select all stored values from database")
	selectByID           = flag.String("id", "", "select by id")
)

func main() {
	flag.Parse()

	var driverName, dataSourceName string
	if *usingMysql {
		driverName, dataSourceName = "mysql", fmt.Sprintf("test:test@tcp(localhost:%d)/test?loc=Local&charset=utf8mb4&timeout=90s&collation=utf8mb4_unicode_ci", *dbport)
	} else {
		driverName, dataSourceName = "postgres", fmt.Sprintf("sslmode=disable dbname=test user=test password=test host=127.0.0.1 port=%d", *dbport)
	}

	db, err := sql.Open(driverName, dataSourceName)
	logFatal(err)
	defer db.Close()

	err = db.Ping()
	logFatal(err)

	if *createTable {
		_, err := db.Exec("DROP TABLE IF EXISTS test_table")
		logFatal(err)

		if *usingMysql {
			_, err = db.Exec("CREATE TABLE IF NOT EXISTS test_table(id int NOT NULL AUTO_INCREMENT PRIMARY KEY, username VARBINARY(2000), password VARBINARY(2000), email VARBINARY(2000))")
			logFatal(err)
		} else {
			_, err = db.Exec("DROP SEQUENCE IF EXISTS test_table_seq")
			logFatal(err)
			_, err = db.Exec("CREATE SEQUENCE test_table_seq START 1")
			logFatal(err)
			_, err = db.Exec("CREATE TABLE IF NOT EXISTS test_table(id INTEGER PRIMARY KEY DEFAULT nextval('test_table_seq'), username BYTEA, password BYTEA, email BYTEA)")
			logFatal(err)
		}
		log.Println("Table has been successfully created")
	}

	if *dropTable {
		_, err := db.Exec("DROP TABLE IF EXISTS test_table")
		logFatal(err)
		if !*usingMysql {
			_, err = db.Exec("DROP SEQUENCE IF EXISTS test_table_seq")
		}
		logFatal(err)
		log.Println("Table has been successfully dropped")
	}

	if *insertRandomValues > 0 {
		if *insertRandomValues > MAXRANDOM {
			log.Fatal("Too much to insert. Use value from range [1 .. " + fmt.Sprint(MAXRANDOM) + "]")
		}

		emails, err := loadFile("demo/testdata/emails")
		logFatal(err)
		passwords, err := loadFile("demo/testdata/passwords")
		logFatal(err)
		usernames, err := loadFile("demo/testdata/usernames")
		logFatal(err)

		s1 := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s1)

		for i := 0; i < *insertRandomValues; i++ {
			userName := getRandomInput(r, usernames)
			password := getRandomInput(r, passwords)
			email := getRandomInput(r, emails)
			_, err := db.Exec(`insert into test_table(username, password, email) values ('` + userName + `', '` + password + `', '` + email + `')`)
			logFatal(err)
		}

		log.Println("Insert has been successful")
	}

	if *poisonRecordToInsert != "" {
		value, err := base64.StdEncoding.DecodeString(*poisonRecordToInsert)
		logFatal(err)

		_, err = db.Exec(`insert into test_table(username, password, email) values ('poison_record', '\x` + hex.EncodeToString(value) + `', '\x` + hex.EncodeToString(value) + `')`)
		logFatal(err)
		log.Println("Poison record insert has been successful")
	}

	if *selectAllFromTable > 0 || *selectByID != "" {
		q := "select * from test_table"
		if *selectByID != "" {
			q += " where id=" + *selectByID
		}
		if *selectAllFromTable > 0 {
			q += " limit " + strconv.Itoa(*selectAllFromTable)
		}
		rows, err := db.Query(q)
		logFatal(err)
		type Row struct {
			id       int
			username []byte
			password []byte
			email    []byte
		}
		for rows.Next() {
			var r Row
			err := rows.Scan(&r.id, &r.username, &r.password, &r.email)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%d\t%s\t%s\t%s\n", r.id, tryString(r.username), tryString(r.password), tryString(r.email))
		}
		rows.Close()

		log.Println("Select has been successful")
	}
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const MAXRANDOM = 25

func getRandomInput(r *rand.Rand, from []string) string {
	return from[r.Intn(MAXRANDOM)]
}

func loadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

// tryString tries to convert byte slice into valid utf8 string
// if conversion is unsuccessfull, the hex string is returned with leading '\x'
func tryString(slice []byte) string {
	if utf8.Valid(slice) {
		return string(slice)
	}

	if *usingHex {
		return "0x" + hex.EncodeToString(slice)
	}

	return base64.URLEncoding.EncodeToString(slice)
}
