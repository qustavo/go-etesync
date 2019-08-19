package cache

import (
	"fmt"

	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/store"
)

type Cache struct {
	store store.Store
	api   api.Client
}

func New(s store.Store, c api.Client) *Cache {
	return &Cache{store: s, api: c}
}

// Sync syncs all the available journals
func (c *Cache) Sync() error {
	js, err := c.Journals()
	if err != nil {
		return err
	}

	for _, j := range js {
		if err := c.SyncJournal(j.UID); err != nil {
			return err
		}
	}
	return nil
}

// SyncJournal write to the last entries (using the ?last arg) to the store
func (c *Cache) SyncJournal(uid string) error {
	e, err := c.store.LastEntry(uid)
	if err != nil && err != store.ErrRecordNotFound {
		return err
	}

	var last *string = nil
	if err != store.ErrRecordNotFound {
		last = &e.UID
	}
	entries, err := c.api.JournalEntries(uid, last)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if err := c.store.CreateEntry(uid, e); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) Journals() (api.Journals, error) {
	js, err := c.api.Journals()
	if err != nil {
		return nil, err
	}

	var v2Journals api.Journals
	for _, j := range js {
		if j.Version != 2 {
			fmt.Printf("[skip] Journal:%s (version %d != 2)\n", j.UID, j.Version)
			continue
		}
		v2Journals = append(v2Journals, j)
	}

	return v2Journals, nil
}

func (c *Cache) JournalEntries(uid string) (api.Entries, error) {
	return c.store.GetEntries(uid)
}
