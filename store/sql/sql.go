package sql

import (
	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/store"
	"github.com/jinzhu/gorm"
)

type Entry struct {
	ID         uint   `gorm:"primary_key"`
	JournalUID string `gorm:"index:journal_uid;not null"`
	*api.Entry
}

type Store struct {
	db *gorm.DB
}

func NewStore(driver, dsn string) (*Store, error) {
	db, err := gorm.Open(driver, dsn)
	return &Store{db: db}, err
}

func (s *Store) Migrate() error {
	err := s.db.AutoMigrate(
		&api.Journal{},
		&Entry{},
	).Error
	if err != nil {
		return err
	}

	return nil
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

func (s *Store) CreateEntry(j string, e *api.Entry) error {
	return s.db.Create(
		&Entry{JournalUID: j, Entry: e},
	).Error
}

func (s *Store) GetEntry(j string, uid string) (*api.Entry, error) {
	var e = &api.Entry{}
	if err := s.first("uid = ?", uid, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Store) GetEntries(j string) (api.Entries, error) {
	var entries api.Entries
	db := s.db.Find(&entries, "journal_uid = ?", j)
	if db.RecordNotFound() {
		return nil, nil
	}

	if err := db.Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *Store) LastEntry(j string) (*api.Entry, error) {
	var e Entry
	db := s.db.Last(&e, "journal_uid = ?", j)
	if db.RecordNotFound() {
		return nil, store.ErrRecordNotFound
	}

	if err := db.Error; err != nil {
		return nil, err
	}

	return e.Entry, nil
}
