package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	user      = "postgres"
	password  = "postgres123"
	host      = "localhost"
	port      = 5432
	dbname    = "postgres"
	sslmode   = "disable"
	maxClient = 10
)

func main1() {
	fmt.Println("db client")
	dbInfo := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%v", user, password, host, port, dbname, sslmode)
	open, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer open.Close()
	open.SetMaxOpenConns(10)

	/*emps, names := getDataFromExcel()
	createTableInDB(open)
	addDataToDB(open, emps, names)
	fmt.Println("done")*/

	wg := sync.WaitGroup{}
	wg.Add(maxClient)

	go func() {
		for i := 0; i < maxClient; i++ {
			queryDataFromDB(open)
			wg.Done()
		}
	}()
	wg.Wait()
}

func createTableInDB(db *sql.DB) {
	createTableQuery := `
			CREATE TABLE IF NOT EXISTS TEST_TABLE (
    EEID VARCHAR,
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
    EXIT_DATE VARCHAR,
    PRIMARY KEY (EEID, FULL_NAME)
) PARTITION BY RANGE (FULL_NAME);`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
		return
	}

	createNameTableQuery := `
			CREATE TABLE IF NOT EXISTS TEST_NAME (
    ID VARCHAR,
    NAME1 VARCHAR,
    PRIMARY KEY (ID, NAME1)
) PARTITION BY RANGE (NAME1) ;`

	_, err = db.Exec(createNameTableQuery)
	if err != nil {
		log.Fatal(err)
		return
	}

	pq := `
	CREATE TABLE TEST_TABLE_AD PARTITION OF TEST_TABLE
	FOR VALUES FROM ('A') TO ('E');

	
	CREATE TABLE TEST_TABLE_EH PARTITION OF TEST_TABLE
	FOR VALUES FROM ('E') TO ('I');

	
	CREATE TABLE TEST_TABLE_IL PARTITION OF TEST_TABLE
	FOR VALUES FROM ('I') TO ('M');

	
	CREATE TABLE TEST_TABLE_MP PARTITION OF TEST_TABLE
	FOR VALUES FROM ('M') TO ('Q');

	
	CREATE TABLE TEST_TABLE_QT PARTITION OF TEST_TABLE
	FOR VALUES FROM ('Q') TO ('U');

	
	CREATE TABLE TEST_TABLE_UX PARTITION OF TEST_TABLE
	FOR VALUES FROM ('U') TO ('Y');

	
	CREATE TABLE TEST_TABLE_YZ PARTITION OF TEST_TABLE
	FOR VALUES FROM ('Y') TO ('ZZZ');
`
	_, err = db.Exec(pq)
	if err != nil {
		fmt.Println(err)
		return
	}
	pq1 := `
	CREATE TABLE TEST_NAME_AD PARTITION OF TEST_NAME
	FOR VALUES FROM ('A') TO ('E');

	
	CREATE TABLE TEST_NAME_EH PARTITION OF TEST_NAME
	FOR VALUES FROM ('E') TO ('I');

	
	CREATE TABLE TEST_NAME_IL PARTITION OF TEST_NAME
	FOR VALUES FROM ('I') TO ('M');

	
	CREATE TABLE TEST_NAME_MP PARTITION OF TEST_NAME
	FOR VALUES FROM ('M') TO ('Q');

	
	CREATE TABLE TEST_NAME_QT PARTITION OF TEST_NAME
	FOR VALUES FROM ('Q') TO ('U');

	
	CREATE TABLE TEST_NAME_UX PARTITION OF TEST_NAME
	FOR VALUES FROM ('U') TO ('Y');

	
	CREATE TABLE TEST_NAME_YZ PARTITION OF TEST_NAME
	FOR VALUES FROM ('Y') TO ('ZZZ');
`

	_, err = db.Exec(pq1)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

}

func addDataToDB(db *sql.DB, employees []employee, names []Name) {
	tx, _ := db.Begin()
	if employees != nil {
		stmt, _ := db.Prepare(`INSERT INTO TEST_TABLE (
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

	stmt1, _ := db.Prepare(`INSERT INTO TEST_NAME (ID, NAME1) VALUES ($1, $2)`)

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

type employee struct {
	ID           string `json:"EEID"`
	Name         string `json:"Full Name"`
	Profession   string `json:"Job Title"`
	Department   string `json:"Department"`
	BusinessUnit string `json:"Business Unit"`
	Gender       string `json:"Gender"`
	Ethnicity    string `json:"Ethnicity"`
	Age          int    `json:"Age"`
	HireDate     string `json:"Hire Date"`
	AnnualSalary string `json:"Annual Salary"`
	Bonus        int    `json:"Bonus"`
	Country      string `json:"Country"`
	City         string `json:"City"`
	ExitDate     string `json:"EEIExit DateD"`
}

func getDataFromExcel() ([]employee, []Name) {

	file, err := excelize.OpenFile("/Users/DASDAIS/IdeaProjects/system-design/partitioning/test_data.xlsx")
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	rows, err := file.GetRows("Data")
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	employes := make([]employee, 0)

	for _, r := range rows {
		expectedCols := 14 // If you expect 14 columns
		row := padRow(r, expectedCols)
		var emp employee
		emp.ID = row[0]
		emp.Name = row[1]
		emp.Profession = row[2]
		emp.Department = row[3]
		emp.BusinessUnit = row[4]
		emp.Gender = row[5]
		emp.Ethnicity = row[6]
		emp.Age, _ = strconv.Atoi(row[7])
		emp.HireDate = row[8]
		emp.AnnualSalary = row[9]
		emp.Bonus, _ = strconv.Atoi(row[10])
		emp.Country = row[11]
		emp.City = row[12]
		emp.ExitDate = row[13]
		employes = append(employes, emp)
	}

	rowsName, _ := file.GetRows("Name")
	names := make([]Name, 0)

	for _, r := range rowsName {
		var name Name
		name.ID = r[0]
		name.Name = r[1]
		names = append(names, name)
	}

	return employes, names
}

type Name struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

func padRow(row []string, desiredLen int) []string {
	if len(row) >= desiredLen {
		return row
	}
	padded := make([]string, desiredLen)
	copy(padded, row)
	// Remaining elements are already "" (zero value)
	return padded
}

func queryDataFromDB(db *sql.DB) {

	t1 := time.Now()
	searchPattern := "A"
	p2 := getPartionTestTable(searchPattern)
	p1 := getPartionTestNameTable(searchPattern)

	query := fmt.Sprintf(`SELECT FULL_NAME, ID FROM 
                         %s t INNER JOIN %s n ON t.FULL_NAME = n.NAME1 WHERE t.FULL_NAME LIKE $1`, p1, p2)

	// Add wildcards for contains search

	rows, err := db.Query(query, "%"+"a"+"%")

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

func getPartionTestNameTable(name string) interface{} {
	if len(name) == 0 {
		return "TEST_TABLE"
	}
	first := string([]rune(name)[0])
	switch {
	case first >= "A" && first <= "D":
		return "TEST_TABLE_AD"
	case first >= "E" && first <= "H":
		return "TEST_TABLE_EH"
	case first >= "I" && first <= "L":
		return "TEST_TABLE_IL"
	case first >= "M" && first <= "P":
		return "TEST_TABLE_MP"
	case first >= "Q" && first <= "T":
		return "TEST_TABLE_QT"
	case first >= "U" && first <= "X":
		return "TEST_TABLE_UX"
	default:
		return "TEST_TABLE_YZ"
	}
}

func getPartionTestTable(name string) interface{} {
	if len(name) == 0 {
		return "TEST_NAME"
	}
	first := string([]rune(name)[0])
	switch {
	case first >= "A" && first <= "D":
		return "TEST_NAME_AD"
	case first >= "E" && first <= "H":
		return "TEST_NAME_EH"
	case first >= "I" && first <= "L":
		return "TEST_NAME_IL"
	case first >= "M" && first <= "P":
		return "TEST_NAME_MP"
	case first >= "Q" && first <= "T":
		return "TEST_NAME_QT"
	case first >= "U" && first <= "X":
		return "TEST_NAME_UX"
	default:
		return "TEST_NAME_YZ"
	}
}
