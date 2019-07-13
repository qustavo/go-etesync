package manager

import (
	"testing"

	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/store/sql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/require"
)

func TestAddEntries(t *testing.T) {
	s, err := sql.NewStore("sqlite3", ":memory:")
	require.NoError(t, err)

	pass := []byte("pass")
	j := &api.Journal{UID: "abcd", Owner: "someones@email"}

	contents := []*api.EntryContent{
		{Action: "ADD", Content: "abcd"},
	}

	entries := make(api.Entries, len(contents))
	for i, c := range contents {
		e := api.NewEntry(j)
		err := e.SetContent(c, pass)
		require.NoError(t, err)
		entries[i] = e
	}

	require.NoError(t, AddEntries(s, entries))
}
