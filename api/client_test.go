package api

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	u := os.Getenv("ETESYNC_TEST_USERNAME")
	p := os.Getenv("ETESYNC_TEST_PASSWORD")
	k := os.Getenv("ETESYNC_TEST_ENCRYPTION")
	if u == "" || p == "" || k == "" {
		t.Skipf("ETESYNC_TEST_{USERNAME || PASSWORD || ENCRYPTION } not declared")
	}

	c, err := NewClient(u, p)
	require.NoError(t, err)

	js, err := c.Journals()
	require.NoError(t, err)

	key, err := DeriveKey(u, []byte(k))
	require.NoError(t, err)

	for _, j := range js {
		dec, err := j.GetContent(key)
		require.NoError(t, err)
		log.Printf("journal: %v", dec)

		es, err := c.JournalEntries(j.UID)
		require.NoError(t, err)

		for _, e := range es {
			dec, err := e.GetContent(key)
			require.NoError(t, err)
			log.Printf("entry: %v", dec)
		}
	}
}
