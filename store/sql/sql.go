package sql

import (
	"github.com/gchaincl/go-etesync/store"
	"github.com/jinzhu/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(driver, dsn string) (*Store, error) {
	db, err := gorm.Open(driver, dsn)
	return &Store{db: db.Debug()}, err
}

func (s *Store) Migrate() error {
	s.db.AutoMigrate(
		&store.Contact{},
		&store.Event{},
	)

	err := s.db.
		Model(&store.Event{}).AddIndex("idx_event_uid", "uid").
		Model(&store.Contact{}).AddIndex("idx_contact_uid", "uid").
		Error

	return err
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) first(where string, val interface{}, model interface{}) error {
	db := s.db.Where(where, val).First(model)
	if db.RecordNotFound() {
		return store.ErrRecordNotFound
	}

	if err := db.Error; err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateJournal() {
	panic("not implemented")
}

func (s *Store) GetJournal() {
	panic("not implemented")
}

func (s *Store) CreateEntry() {
	panic("not implemented")
}

func (s *Store) GetEntry() {
	panic("not implemented")
}

func (s *Store) CreateContact(c *store.Contact) error {
	return s.db.Create(c).Error
}

func (s *Store) GetContact(uid string) (*store.Contact, error) {
	var c = &store.Contact{}
	if err := s.first("uid = ?", uid, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Store) CreateEvent(e *store.Event) error {
	return s.db.Create(e).Error
}

func (s *Store) GetEvent(uid string) (*store.Event, error) {
	var e = &store.Event{}
	if err := s.first("uid = ?", uid, e); err != nil {
		return nil, err
	}
	return e, nil
}
