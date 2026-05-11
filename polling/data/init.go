package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
	"time"
)

const (
	host     = "0.0.0.0"
	port     = 5432
	user     = "postgres"
	password = "postgres123"
	dbname   = "postgres"
)

var (
	once   *sync.Once
	testDB TestDB
)

type TestDB struct {
	db *sql.DB
}

func Initialize() TestDB {
	once = &sync.Once{}
	once.Do(func() {
		pgInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		open, err := sql.Open("postgres", pgInfo)
		fmt.Println(err)
		testDB = TestDB{db: open}
	})
	return testDB
}

func InsertData() {

	createQuery := "CREATE TABLE IF NOT EXISTS USERS(id INTEGER PRIMARY KEY, STATUS VARCHAR(20));"

	CreateInsertTable(createQuery)

	insertQuery := fmt.Sprintf("INSERT INTO USERS (id, STATUS) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET STATUS = $2")
	CreateInsertTable(insertQuery, 1, "CREATED")

	time.Sleep(2 * time.Minute)

	insertQuery = fmt.Sprintf("INSERT INTO USERS(id, STATUS) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET STATUS = %s;", "IN PROGRESS")
	CreateInsertTable(insertQuery, 1, "CREATED")

	time.Sleep(2 * time.Minute)

	insertQuery = fmt.Sprintf("INSERT INTO USERS(id, STATUS) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET STATUS = %s;", "DONE")
	CreateInsertTable(insertQuery, 1, "CREATED")

}
