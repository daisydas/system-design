package main

import (
	"database/sql"
	"fmt"
	"sync"
)

type flightSeatHandler struct {
	flightDB *sql.DB
}

func (fsh *flightSeatHandler) AllocateSeats(pnr string, wg *sync.WaitGroup) {

	defer wg.Done()
	tx, err := fsh.flightDB.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}
	var seatID string
	row := tx.QueryRow("SELECT SEATID FROM PLANE WHERE PNR IS NULL LIMIT 1 FOR UPDATE SKIP LOCKED")
	_ = row.Scan(&seatID)

	_, _ = tx.Exec("UPDATE PLANE SET PNR = $1 WHERE SEATID = $2", pnr, seatID)

	tx.Commit()
}
