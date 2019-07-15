package storetest

import (
	"fmt"
	"testing"

	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Factory func(t *testing.T) (store.Store, func())

func TestSuite(t *testing.T, f Factory) {
	tests := []struct {
		name string
		run  func(*testing.T, store.Store)
	}{
		{"Entry/Create", TestEntryCreate},
		{"Entry/NotFound", TestEntryNotFound},
		{"Entry/Last", TestLastEntry},
		{"Entry/GetEntries", TestEntryGetEntries},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, cleanup := f(t)
			defer cleanup()

			test.run(t, s)
		})
	}
}

func TestEntryCreate(t *testing.T, s store.Store) {
	e := &api.Entry{UID: "abcd", Content: "data"}

	err := s.CreateEntry("parent", e)
	require.NoError(t, err)

	found, err := s.GetEntry("parent", e.UID)
	require.NoError(t, err)

	assert.Equal(t, e, found)
}

func TestEntryNotFound(t *testing.T, s store.Store) {
	notFound, err := s.GetEntry("parent", "abcd")
	require.Error(t, err)
	assert.Equal(t, store.ErrRecordNotFound, err)
	assert.Nil(t, notFound)
}

func TestLastEntry(t *testing.T, s store.Store) {
	entries := api.Entries{
		&api.Entry{UID: "01"},
		&api.Entry{UID: "02"},
		&api.Entry{UID: "03"},
		&api.Entry{UID: "04"},
		&api.Entry{UID: "05"},
		&api.Entry{UID: "06"},
	}
	for _, e := range entries {
		err := s.CreateEntry("parent", e)
		require.NoError(t, err)

		found, err := s.LastEntry("parent")
		require.NoError(t, err)
		assert.Equal(t, e.UID, found.UID)
	}
}

func TestEntryGetEntries(t *testing.T, s store.Store) {
	total := 100
	for i := 0; i < total; i++ {
		e := &api.Entry{UID: fmt.Sprintf("uid%d", i)}

		require.NoError(t, s.CreateEntry("uid-a", e))
		require.NoError(t, s.CreateEntry("uid-b", e))

	}

	entries, err := s.GetEntries("uid-a")
	require.NoError(t, err)

	assert.Len(t, entries, total)

	t.Run("empty entries", func(t *testing.T) {
		entries, err := s.GetEntries("uid-xxx")
		require.NoError(t, err)
		assert.Len(t, entries, 0)
	})
}
