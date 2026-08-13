package connectionpool

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5420
	user     = "postgres"
	password = "postgres123"
	dbname   = "postgres"
)

type Conn struct {
	db  *sql.DB
	dbs []*sql.DB
}

func getDBConnection(id int) int {
	if id%2 == 0 {
		return 0
	}
	return 1
}

func newConnection() *Conn {
	pgInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, _ := sql.Open("postgres", pgInfo)
	return &Conn{db: db}
}

func newConnectionBasedOnID() *Conn {

	dbs := make([]*sql.DB, 0)
	pgInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db1, _ := sql.Open("postgres", pgInfo)
	dbs = append(dbs, db1)
	db2, _ := sql.Open("postgres", pgInfo)
	dbs = append(dbs, db2)

	return &Conn{dbs: dbs}
}

func main() {
	conn := newConnectionBasedOnID()
	ge := gin.Default()
	ge.POST("/heartbeats", func(ctx *gin.Context) {
		idQ := ctx.Query("id")
		id, _ := strconv.Atoi(idQ)
		insertStm := getInsertQuery(id)
		_, err := conn.dbs[getDBConnection(id)].Exec(insertStm)
		if err != nil {
			fmt.Println(err)
			return
		}
		ctx.JSON(200, gin.H{"message": "ok"})
	})

	ge.GET("/heartbeats/status", func(ctx *gin.Context) {
		x, y := getQuery(ctx.Query("ids"))
		rows, err := conn.db.Query(x, y)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp := make([]map[string]bool, 0)
		var userID, lastHB int
		for rows.Next() {
			_ = rows.Scan(&userID, &lastHB)
			resp = append(resp, map[string]bool{fmt.Sprintf("%d", userID): lastHB > int(time.Now().Unix())-30})
		}

		ctx.JSON(200, gin.H{"message": resp})
	})
	_ = ge.Run(":8080")

}

func getInsertQuery(userID int) string {
	tt := time.Now().Unix()
	return fmt.Sprintf("INSERT INTO users(ID, EPOC_TIME) VALUES (%d, %v) ON CONFLICT (ID) DO UPDATE SET ID = %d, EPOC_TIME = %v", userID, tt, userID, tt)
}

func getQuery(userIDs string) (string, []int) {
	ids := strings.Split(userIDs, ",")
	idsnum := make([]int, 0)
	for _, id := range ids {
		v, _ := strconv.Atoi(id)
		idsnum = append(idsnum, v)
	}
	placeholders := strings.Repeat("?,", len(ids)-1) + "?"

	return fmt.Sprintf("SELECT EPOC_TIME FROM users WHERE ID (%s)", placeholders), idsnum
}
