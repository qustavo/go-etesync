package api

import "github.com/gchaincl/go-etesync/crypto"

func DeriveKey(email string, password []byte) ([]byte, error) {
	return crypto.DeriveKey(password, []byte(email))
}
