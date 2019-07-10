package sql

import (
	"testing"

	"github.com/gchaincl/go-etesync/store"
	"github.com/gchaincl/go-etesync/store/storetest"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestSQL(t *testing.T) {
	f := func(t *testing.T) (store.Store, func()) {
		s, err := NewStore("sqlite3", ":memory:")
		require.NoError(t, err)
		require.NoError(t, s.Migrate())

		return s, func() { s.Close() }

	}
	storetest.TestSuite(t, f)
}
