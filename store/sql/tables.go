package sql

import "github.com/gchaincl/go-etesync/api"

type Entry struct {
	ID         uint   `gorm:"primary_key"`
	JournalUID string `gorm:"index:journal_uid;not null"`
	*api.Entry
}
