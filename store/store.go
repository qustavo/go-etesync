package store

import (
	"errors"

	"github.com/gchaincl/go-etesync/api"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Store interface {
	CreateJournal(*api.Journal) error
	GetJournal()

	CreateEntry(*api.Entry) error
	GetEntry(uid string) (*api.Entry, error)
}
