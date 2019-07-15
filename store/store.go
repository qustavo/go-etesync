package store

import (
	"errors"

	"github.com/gchaincl/go-etesync/api"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Store interface {
	CreateEntry(journalUID string, entry *api.Entry) error
	GetEntries(journalUID string) (api.Entries, error)
	GetEntry(journalUID string, entryUID string) (*api.Entry, error)
	LastEntry(journalUID string) (*api.Entry, error)
}
