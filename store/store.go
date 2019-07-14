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

	CreateEntry(journalUID string, entry *api.Entry) error
	GetEntry(journalUID string, entryUID string) (*api.Entry, error)
	LastEntry(journalUID string) (*api.Entry, error)
}
