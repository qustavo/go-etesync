package storetest

import (
	"testing"
	"time"

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
		{"Entity/CreateContact", TestEntityCreateContact},
		{"Entity/ContactNotFound", TestEntityContactNotFound},
		{"Entity/CreateEvent", TestEntityCreateEvent},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, cleanup := f(t)
			defer cleanup()

			test.run(t, s)
		})
	}
}

func TestEntityCreateContact(t *testing.T, s store.Store) {
	c := &store.Contact{
		UID: "abcd", Name: "alice", Phone: "123123123",
	}

	err := s.CreateContact(c)
	require.NoError(t, err)

	found, err := s.GetContact(c.UID)
	require.NoError(t, err)

	assert.Equal(t, c, found)
}

func TestEntityContactNotFound(t *testing.T, s store.Store) {
	notFound, err := s.GetContact("abcd")
	require.Error(t, err)
	assert.Equal(t, store.ErrRecordNotFound, err)
	assert.Nil(t, notFound)
}

func TestEntityCreateEvent(t *testing.T, s store.Store) {
	e := &store.Event{
		UID: "abcd", Summary: "test", Date: time.Date(2009, 1, 3, 0, 0, 0, 0, time.UTC),
	}

	err := s.CreateEvent(e)
	require.NoError(t, err)

	found, err := s.GetEvent(e.UID)
	require.NoError(t, err)

	assert.Equal(t, e, found)
}
