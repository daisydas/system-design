package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	user     = "postgres"
	password = "postgres123"
	host     = "localhost"
	port     = "5432"
	dbname   = "postgres"
)

type Operation struct {
	db *sql.DB
}

func (ops *Operation) GetFromKVStore(key string) string {
	var value string
	err := ops.db.QueryRow("SELECT VAL FROM KEY_VALUE WHERE EXPIRED_AT > $1 AND KEY = $2", time.Now().Unix(), key).Scan(&value)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return value
}

func (ops *Operation) PutToKVStore(key, value string, ttlInSecs int) error {
	expiredAt := time.Now().Unix() + int64(ttlInSecs)
	exec, err := ops.db.Exec(`
    INSERT INTO KEY_VALUE (KEY, VALUE, EXPIRED_AT) 
    VALUES ($1, $2, $3)
    ON CONFLICT (KEY) 
    DO UPDATE SET 
        VALUE = EXCLUDED.VALUE,
        EXPIRED_AT = EXCLUDED.EXPIRED_AT`, key, value, expiredAt)
	if err != nil {
		fmt.Println(err)
	}

	if _, err = exec.RowsAffected(); err != nil {
		fmt.Println(err)
	}

	return err

}

func (ops *Operation) DeleteFromKVStore(key string) error {
	exec, err := ops.db.Exec("UPDATE KEY_VALUE SET EXPIRED_AT = -1 WHERE KEY = $1 AND EXPIRED_AT >= $2", key, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
	}

	if _, err = exec.RowsAffected(); err != nil {
		fmt.Println(err)
	}

	return err
}

func (ops *Operation) CreateTable() {

	tx, _ := ops.db.Begin()

	_, err := tx.Exec(`
    CREATE TABLE IF NOT EXISTS KEY_VALUE (
        KEY VARCHAR(20) PRIMARY KEY, 
        VALUE VARCHAR(20), 
        EXPIRED_AT INTEGER
    ) PARTITION BY HASH (KEY)
`)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}

	_, err = tx.Exec(`CREATE TABLE KEY_VALUE_P0 PARTITION OF KEY_VALUE FOR VALUES WITH(MODULUS 2, REMAINDER 0)`)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`CREATE TABLE KEY_VALUE_P1 PARTITION OF KEY_VALUE FOR VALUES WITH(MODULUS 2, REMAINDER 1)`)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	return
}

func (ops *Operation) DropTable() {
	_, err := ops.db.Exec("DROP TABLE IF EXISTS KEY_VALUE")
	if err != nil {
		fmt.Println(err)
	}
	return
}

func NewOperation() *Operation {
	open, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port = %v sslmode=disable", user, password, dbname, host, port))
	if err != nil {
		return nil
	}
	return &Operation{open}
}
