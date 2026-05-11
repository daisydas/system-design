package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"system-design/polling/data"
	"system-design/polling/handler"
)

func main() {

	ge := gin.Default()
	db := data.Initialize()
	go data.InsertData()

	sph := handler.NewShortPollHandler(db)
	lph := handler.NewLongPollHandler(db)

	ge.GET("/poll/short", func(c *gin.Context) {
		idStr := c.Query("id")
		id, _ := strconv.Atoi(idStr)
		resp := sph.PollStatus(id)
		c.JSON(http.StatusOK, resp)
	})

	ge.GET("/poll/long", func(c *gin.Context) {
		idStr := c.Query("id")
		prev := c.Query("prev")
		id, _ := strconv.Atoi(idStr)
		resp := lph.PollStatus(prev, id)
		c.JSON(http.StatusOK, resp)
	})
	_ = ge.Run(":8097")

}
