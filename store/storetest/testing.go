package storetest

import (
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

	err := s.CreateEntry(e)
	require.NoError(t, err)

	found, err := s.GetEntry(e.UID)
	require.NoError(t, err)

	assert.Equal(t, e, found)
}

func TestEntryNotFound(t *testing.T, s store.Store) {
	notFound, err := s.GetEntry("abcd")
	require.Error(t, err)
	assert.Equal(t, store.ErrRecordNotFound, err)
	assert.Nil(t, notFound)
}
