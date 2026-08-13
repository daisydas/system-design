package service

import "os"

type authorizationFunc func() bool

func (s authorizationFunc) display() {

}

type DB struct {
	authorizationFunc authorizationFunc
}

func (d *DB) IsAuthorized() bool {
	return d.authorizationFunc()
}

func NewDb() *DB {
	return &DB{authorizationFunc: authorize}
}

func authorize() bool {

	ts := authorizationFunc(func() bool { return true })
	ts.display()

	user := os.Args[1]
	if user == "admin" {
		return true
	}
	return false
}
