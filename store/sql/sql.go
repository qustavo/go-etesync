package sql

import (
	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/store"
	"github.com/jinzhu/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(driver, dsn string) (*Store, error) {
	db, err := gorm.Open(driver, dsn)
	return &Store{db: db}, err
}

func (s *Store) Migrate() error {
	return s.db.AutoMigrate(
		&api.Entry{},
	).Error
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

func (s *Store) CreateJournal(j *api.Journal) error {
	panic("not implemented")
}

func (s *Store) GetJournal() {
	panic("not implemented")
}

func (s *Store) CreateEntry(e *api.Entry) error {
	return s.db.Create(e).Error
}

func (s *Store) GetEntry(uid string) (*api.Entry, error) {
	var e = &api.Entry{}
	if err := s.first("uid = ?", uid, e); err != nil {
		return nil, err
	}
	return e, nil
}
