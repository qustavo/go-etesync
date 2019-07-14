package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCipher(t *testing.T) {
	m := New([]byte("salt"), []byte("key"))
	plaintext := []byte("0000000000000000X")
	enc, err := m.Encrypt(plaintext)
	require.NoError(t, err)

	dec, err := m.Decrypt(enc)
	require.NoError(t, err)

	assert.Equal(t, plaintext, dec)
}
