package main

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNoConnectionPool(t *testing.T) {
	tt := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(5000)
	for i := 0; i < 5000; i++ {
		go func() {
			pgInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
			db, _ := sql.Open("postgres", pgInfo)

			_, _ = db.Exec("SELECT SLEEP(0.01)")
			_ = db.Close()
			defer wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(tt))

}

func TestConnectionPool(t *testing.T) {
	tt := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(5000)

	cp := NewConnectionPool(30)
	for i := 0; i < 5000; i++ {
		go func() {
			con := cp.Get()
			_, _ = con.db.Exec("SELECT SLEEP(0.01)")
			cp.Release(con)
			defer wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Since(tt))
	defer cp.Close()
}
