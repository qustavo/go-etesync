package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryContentEncryption(t *testing.T) {
	jn := &Journal{UID: "abcd", Owner: "some@email"}
	ec := &EntryContent{Action: "ADD", Content: "string"}

	key := []byte("encryption key")
	cipher := crypo.New(jn.UID, key

	en := &Entry{}
	err := en.SetContent(ec, pass)
	require.NoError(t, err)

	newEc, err := en.GetContent(pass)
	require.NoError(t, err)

	assert.Equal(t, ec, newEc)
}
