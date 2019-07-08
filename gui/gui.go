package gui

import (
	"github.com/gchaincl/go-etesync/api"
	"github.com/rivo/tview"
)

func Start(c *api.Client, key []byte) error {
	box := tview.NewBox().SetBorder(true).SetTitle("ETESync")
	return tview.NewApplication().SetRoot(box, true).Run()
}
