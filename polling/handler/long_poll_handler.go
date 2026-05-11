package handler

import (
	"fmt"
	"system-design/polling/data"
	"system-design/polling/model"
	"time"
)

const timeOutInMins = 2

type LongPollHandler struct {
	testDatabase data.TestDB
}

func (lph *LongPollHandler) PollStatus(prev string, id int) model.Resp {
	queryStr := fmt.Sprintf("SELECT status FROM USERS WHERE ID = %d", id)
	t1 := time.Now()
	for {
		status := data.QueryTable(queryStr)
		if prev != status {
			return model.Resp{
				Status: status,
			}
		}
		if time.Since(t1).Minutes() > timeOutInMins {
			fmt.Println("timeout")
			break
		}

	}
	return model.Resp{}
}

func NewLongPollHandler(testDatabase data.TestDB) *LongPollHandler {
	return &LongPollHandler{testDatabase}
}
