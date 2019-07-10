package store

import "time"

type Journal struct {
	Version  int
	UID      string
	Content  interface{}
	Owner    string
	Key      string
	ReadOnly bool
}

type Contact struct {
	UID   string
	Name  string
	Phone string
}

type Event struct {
	UID     string
	Summary string
	Date    time.Time
}
