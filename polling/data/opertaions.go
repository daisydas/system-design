package data

import (
	"fmt"
)

func CreateInsertTable(createInsertTableQuery string, args ...any) {
	var err error
	if args == nil {
		_, err = testDB.db.Exec(createInsertTableQuery)
	} else {
		_, err = testDB.db.Exec(createInsertTableQuery, args[0], args[1])
	}

	if err != nil {
		fmt.Println(err)
		return
	}
}

func QueryTable(query string) string {
	var STATUS string
	row := testDB.db.QueryRow(query)
	err := row.Scan(&STATUS)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return STATUS
}
