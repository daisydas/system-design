package handler

import (
	"fmt"
	"system-design/polling/data"
	"system-design/polling/model"
)

type ShortPollHandler struct{}

func (sph *ShortPollHandler) PollStatus(id int) model.Resp {
	queryString := fmt.Sprintf("SELECT status FROM USERS WHERE ID = %d LIMIT 1", id)
	status := data.QueryTable(queryString)
	return model.Resp{
		Status: status,
	}

}

func NewShortPollHandler(testDatabase data.TestDB) *ShortPollHandler {
	return &ShortPollHandler{}
}
