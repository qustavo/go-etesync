package api

import "encoding/base64"
import "github.com/gchaincl/go-etesync/crypto"

type Journal struct {
	Version  int    `json:"version"`
	UID      string `json:"uid"`
	Content  string `json:"content"`
	Owner    string `json:"owner"`
	Key      string `json:"key"`
	ReadOnly bool   `json:"readOnly"`

	derivedKey []byte
}

func (j *Journal) DerivedKey(password []byte) ([]byte, error) {
	if j.derivedKey != nil {
		return j.derivedKey, nil
	}

	key, err := crypto.DeriveKey(password, []byte(j.Owner))
	if err != nil {
		return nil, err
	}

	j.derivedKey = key
	return key, nil
}

func (j *Journal) GetContent(password []byte) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(j.Content)
	if err != nil {
		return nil, err
	}

	key, err := j.DerivedKey(password)
	if err != nil {
		return nil, err
	}

	return crypto.New([]byte(j.UID), key).Decrypt(content[32:])
}

type Journals []*Journal

type Entry struct {
	UID     string `json:"uid"`
	Content string `json:"content"`
}

func (e *Entry) GetContent(j *Journal, password []byte) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(e.Content)
	if err != nil {
		return nil, err
	}

	key, err := j.DerivedKey(password)
	if err != nil {
		return nil, err
	}

	return crypto.New([]byte(j.UID), key).Decrypt(content[32:])
}

type Entries []Entry
