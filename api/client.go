package api

type Client interface {
	Journals() (Journals, error)
	Journal(uid string) (*Journal, error)
	JournalEntries(uid string, last *string) (Entries, error)
}
