package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
)

const (
	pgHost   = "localhost"
	pgPort   = 5432
	pgDBName = "postgres"
	pgUser   = "postgres"
	pgPass   = "postgres123"
)

func InitializeDB(user []string) *sql.DB {

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pgHost, pgPort, pgUser, pgPass, pgDBName))
	if err != nil {
		//fmt.Println(err)
		return nil
	}

	//db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)

	_, err = db.Exec("CREATE TABLE PLANE (SEATID VARCHAR(20) PRIMARY KEY, PNR VARCHAR(20))")
	if err != nil {
		//fmt.Println(err)
	}
	var seat string
	for i := 1; i <= length/6; i++ {

		for j := 'A'; j <= 'F'; j++ {
			seat = strconv.Itoa(i) + string(j)
			_, err = db.Exec("INSERT INTO PLANE (SEATID, PNR) VALUES ($1, $2)", seat, nil)
			if err != nil {
				//fmt.Println(err)
			}
		}

	}

	db.Exec("CREATE TABLE PLANE_USERS (PNR VARCHAR(20) PRIMARY KEY, NAME VARCHAR(20))")

	for i := 0; i < len(user); i++ {
		db.Exec("INSERT INTO PLANE_USERS (PNR, NAME) VALUES ($1, $2)", "PNR_"+strconv.Itoa(i), user[i])
	}

	return db
}
