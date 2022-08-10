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
	"strings"
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
	insertBind           = flag.Bool("insert_bind", false, "insert values by binding variables")
	insertRandomValues   = flag.Int("insert", 0, "insert N random values into the database")
	dbport               = flag.Int("dbport", 9393, "db port")
	poisonRecordToInsert = flag.String("insert_poison", "", "insert poison record (should be in BASE64 format)")
	selectAllFromTable   = flag.Int("select", 0, "select all stored values from database")
	querySQL             = flag.String("query", "", "query SQL")
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

		var ps *sql.Stmt
		if *insertBind {
			if *usingMysql {
				ps, err = db.Prepare(`insert into test_table(username, password, email) values (?,?,?)`)
			} else {
				ps, err = db.Prepare(`insert into test_table(username, password, email) values ($1, $2, $3)`)
			}
			logFatal(err)
		}

		for i := 0; i < *insertRandomValues; i++ {
			userName := getRandomInput(r, usernames)
			password := getRandomInput(r, passwords)
			email := getRandomInput(r, emails)
			var result sql.Result
			if ps != nil {
				result, err = ps.Exec(userName, password, email)
			} else {
				result, err = db.Exec(`insert into test_table(username, password, email) values ('` + userName + `', '` + password + `', '` + email + `')`)
			}
			logFatal(err)
			lastInsertId, _ := result.LastInsertId()
			rowsAffected, _ := result.RowsAffected()
			log.Printf("lastInsertId: %d, rowsAffected: %d", lastInsertId, rowsAffected)
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

		cols, err := rows.Columns()
		logFatal(err)

		printRowsHeader(cols)

		for rows.Next() {
			var id int
			var username, password, email []byte
			err := rows.Scan(&id, &username, &password, &email)
			logFatal(err)
			fmt.Printf("%d\t%s\t%s\t%s\n", id, tryString(username), tryString(password), tryString(email))
		}
		rows.Close()

		log.Println("Select has been successful")
	}

	*querySQL = strings.TrimSpace(*querySQL)
	upperSQL := strings.ToUpper(*querySQL)
	if upperSQL != "" {
		if strings.HasPrefix(upperSQL, "SELECT") {
			rows, err := db.Query(*querySQL)
			logFatal(err)
			scanRows(rows)
			rows.Close()
		} else if strings.HasPrefix(upperSQL, "INSERT") || strings.HasPrefix(upperSQL, "UPDATE") {

		}
	}
}

func printRowsHeader(cols []string) {
	for i, r := range cols {
		if i > 0 {
			fmt.Print("\t")
		}
		fmt.Printf(r)
	}
	fmt.Println()

	for i := range cols {
		if i > 0 {
			fmt.Print("\t")
		}
		fmt.Printf("---")
	}
	fmt.Println()
}

func scanRows(rows *sql.Rows) {
	cols, err := rows.Columns()
	logFatal(err)

	printRowsHeader(cols)

	// Result is your slice string.
	result := make([]sql.NullString, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range result {
		dest[i] = &result[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err := rows.Scan(dest...)
		logFatal(err)

		for i, r := range result {
			if i > 0 {
				fmt.Print("\t")
			}
			if r.Valid {
				fmt.Printf(tryString([]byte(r.String)))
			} else {
				fmt.Printf("<NULL>")
			}
		}
		fmt.Println()
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
