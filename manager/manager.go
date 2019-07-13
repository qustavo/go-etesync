package manager

import (
	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/store"
)

func AddEntries(s store.Store, entries api.Entries) error {
	for _, e := range entries {
		_ = e
	}
	return nil
}
