package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"sync"
	"time"
)

func main() {
	fmt.Println("db client no")
	dbInfo := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%v", user, password, host, port, dbname, sslmode)
	open, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer open.Close()
	open.SetMaxOpenConns(10)

	/*emp, names := getDataFromExcel()
	createTableInDB1(open)
	addDataToDB1(open, emp, names)
	fmt.Println("done")*/

	wg := sync.WaitGroup{}
	wg.Add(maxClient)

	go func() {
		for i := 0; i < maxClient; i++ {
			queryDataFromDB1(open)
			wg.Done()
		}
	}()
	wg.Wait()
}

func createTableInDB1(db *sql.DB) {
	createTableQuery := `
			CREATE TABLE IF NOT EXISTS TEST_TABLE1 (
    EEID VARCHAR PRIMARY KEY,
    FULL_NAME VARCHAR,
    JOB_TITLE VARCHAR,
    DEPARTMENT VARCHAR,
    BUSINESS_UNIT VARCHAR,
    GENDER VARCHAR,
    ETHNICITY VARCHAR,
    AGE INTEGER,
    HIRE_DATE VARCHAR,
    ANNUAL_SALARY VARCHAR,
    BONUS INTEGER,
    COUNTRY VARCHAR,
    CITY VARCHAR,
    EXIT_DATE VARCHAR
);`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
		return
	}

	createNameTableQuery := `
			CREATE TABLE IF NOT EXISTS TEST_NAME1 (
    ID VARCHAR PRIMARY KEY,
    NAME1 VARCHAR
);`

	_, err = db.Exec(createNameTableQuery)
	if err != nil {
		log.Fatal(err)
		return
	}

}

func addDataToDB1(db *sql.DB, employees []employee, names []Name) {
	tx, _ := db.Begin()
	if employees != nil {
		stmt, _ := db.Prepare(`INSERT INTO TEST_TABLE1 (
        EEID, FULL_NAME, JOB_TITLE, DEPARTMENT, BUSINESS_UNIT, 
        GENDER, ETHNICITY, AGE, HIRE_DATE, ANNUAL_SALARY, 
        BONUS, COUNTRY, CITY, EXIT_DATE
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`)

		defer stmt.Close()

		for i, emp := range employees {
			_, err := stmt.Exec(emp.ID, emp.Name, emp.Profession, emp.Department, emp.BusinessUnit,
				emp.Gender, emp.Ethnicity, emp.Age, emp.HireDate, emp.AnnualSalary, emp.Bonus,
				emp.Country, emp.City, emp.ExitDate)
			i++
			if err != nil {
				fmt.Println(emp.Name)
				fmt.Println(i)
				fmt.Println(err)
				return
			}

		}
		defer func() {
			err := stmt.Close()
			if err != nil {
				log.Fatal(err)
				tx.Rollback()
				return
			}

		}()

	}

	stmt1, _ := db.Prepare(`INSERT INTO TEST_NAME1 (ID, NAME1) VALUES ($1, $2)`)

	for i, name := range names {
		_, err := stmt1.Exec(name.ID, name.Name)
		i++
		if err != nil {
			fmt.Println(name.Name)
			fmt.Println(i)
			fmt.Println(err)
			return
		}

	}
	tx.Commit()

}

func queryDataFromDB1(db *sql.DB) {

	t1 := time.Now()
	query := fmt.Sprintf(`SELECT FULL_NAME, ID FROM 
                         %s t INNER JOIN %s n ON t.FULL_NAME = n.NAME1 WHERE t.FULL_NAME LIKE $1`, "TEST_TABLE1", "TEST_NAME1")
	// Add wildcards for contains search
	searchPattern := "%" + "a" + "%"
	rows, err := db.Query(query, searchPattern)

	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var userName, id string
		if err = rows.Scan(&userName, &id); err != nil {
			log.Fatal(err)
			return
		}
	}
	fmt.Printf("**********************%v\n", time.Since(t1))
}
