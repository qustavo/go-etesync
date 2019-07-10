package store

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
)

type JournalStore interface {
	CreateJournal(*Journal) error
	GetJournal()

	CreateEntry()
	GetEntry()
}

type EntityStore interface {
	CreateContact(*Contact) error
	GetContact(uid string) (*Contact, error)

	CreateEvent(*Event) error
	GetEvent(uid string) (*Event, error)
}

type Store interface {
	JournalStore
	EntityStore
}
