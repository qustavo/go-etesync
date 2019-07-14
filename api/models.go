package api

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gchaincl/go-etesync/crypto"
)

type Journal struct {
	Version  int    `json:"version"`
	UID      string `json:"uid"`
	Content  string `json:"content"`
	Owner    string `json:"owner"`
	Key      string `json:"key"`
	ReadOnly bool   `json:"readOnly"`
}

type JournalType string

const (
	JournalCalendar    JournalType = "CALENDAR"
	JournalAddressBook             = "ADDRESS_BOOK"
	JournalTasks                   = "TASKS"
)

type JournalContent struct {
	Type        JournalType `json:"type"`
	Version     int         `json:"version"`
	Selected    bool        `json:"selected"`
	DisplayName string      `json:"displayName"`
	Color       int         `json:"color"`
}

func (j *Journal) GetContent(cipher *crypto.Cipher) (*JournalContent, error) {
	content, err := base64.StdEncoding.DecodeString(j.Content)
	if err != nil {
		return nil, err
	}

	data, err := cipher.Decrypt(content[32:])
	if err != nil {
		return nil, err
	}

	jc := &JournalContent{}
	if err := json.Unmarshal(data, jc); err != nil {
		return nil, err
	}

	return jc, nil
}

type Journals []*Journal

type Entry struct {
	UID     string `json:"uid"`
	Content string `json:"content"`
}

func (e *Entry) GetContent(cipher *crypto.Cipher) (*EntryContent, error) {
	content, err := base64.StdEncoding.DecodeString(e.Content)
	if err != nil {
		return nil, err
	}

	data, err := cipher.Decrypt(content)
	if err != nil {
		return nil, err
	}

	ec := &EntryContent{}
	if err := json.Unmarshal(data, ec); err != nil {
		return nil, err
	}

	return ec, nil
}

func (e *Entry) SetContent(c *EntryContent, cipher *crypto.Cipher) error {
	json, err := json.Marshal(c)
	if err != nil {
		return err
	}

	data, err := cipher.Encrypt(json)
	if err != nil {
		return err
	}

	e.Content = base64.StdEncoding.EncodeToString(data)
	return nil
}

type Entries []*Entry

type EntryContent struct {
	Action  string
	Content string
}
