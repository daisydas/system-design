package service

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_ToDoWrite(t *testing.T) {
	cases := []struct {
		name    string
		text    string
		input   string
		db      DB
		output  string
		isPanic bool
	}{
		{
			name:    "",
			text:    "abc",
			input:   "xyz",
			db:      DB{authorizationFunc: func() bool { return true }},
			output:  "xyz",
			isPanic: false,
		},
		{
			name:    "",
			text:    "abc",
			input:   "xyz",
			db:      DB{authorizationFunc: func() bool { return false }},
			output:  "abc",
			isPanic: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			todo := Todo{
				Text: c.text,
				db:   c.db,
			}

			todo.Write(c.input)
			if c.isPanic {

			} else {
				assert.Equal(t, c.output, todo.Text)
			}
		})
	}
}
