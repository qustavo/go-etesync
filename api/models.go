package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"

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

func (j *Journal) GetContent(key []byte) (*JournalContent, error) {
	content, err := base64.StdEncoding.DecodeString(j.Content)
	if err != nil {
		return nil, err
	}

	data, err := crypto.New([]byte(j.UID), key).Decrypt(content[32:])
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
	journal *Journal
	UID     string `json:"uid"`
	Content string `json:"content"`
}

func NewEntry(j *Journal) *Entry { return &Entry{journal: j} }

func (e *Entry) Journal() *Journal { return e.journal }

func (e *Entry) GetContent(key []byte) (*EntryContent, error) {
	if e.journal == nil {
		return nil, errors.New(".Journal can't be nil")
	}

	content, err := base64.StdEncoding.DecodeString(e.Content)
	if err != nil {
		return nil, err
	}

	data, err := crypto.New([]byte(e.journal.UID), key).Decrypt(content)
	if err != nil {
		return nil, err
	}

	ec := &EntryContent{}
	if err := json.Unmarshal(data, ec); err != nil {
		return nil, err
	}

	return ec, nil
}

func (e *Entry) SetContent(c *EntryContent, key []byte) error {
	json, err := json.Marshal(c)
	if err != nil {
		return err
	}

	data, err := crypto.New([]byte(e.journal.UID), key).Encrypt(json)
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
