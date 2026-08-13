package service

type Todo struct {
	Text string
	db   DB
}

func (t *Todo) Write(s string) {
	if t.db.IsAuthorized() {
		t.Text = s
	} else {
		panic("user not authorized to write")
	}
}

func (t *Todo) Append(s string) {
	if t.db.IsAuthorized() {
		t.Text += s
	} else {
		panic("user not authorized to write")
	}
}
