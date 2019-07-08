package gui

import (
	"log"
	"time"

	"github.com/gchaincl/go-etesync/api"
	"github.com/laurent22/ical-go"
	"github.com/rivo/tview"
)

func newJournalTable(c *api.Client, key []byte, fn func(*api.Journal)) (*tview.Table, error) {
	js, err := c.Journals()
	if err != nil {
		return nil, err
	}

	t := tview.NewTable().SetSelectable(true, false)
	t.SetTitle("Journals").SetBorder(true)

	uids := make([]*api.Journal, len(js))
	for i, j := range js {
		content, err := j.GetContent(key)
		if err != nil {
			return nil, err
		}
		uids[i] = j

		var icon string

		switch content.Type {
		case api.JournalCalendar:
			icon = "📅"
		case api.JournalAddressBook:
			icon = "🙎"
		case api.JournalTasks:
			icon = "🗒"
		}
		t.SetCell(i, 0, tview.NewTableCell(icon+" "+content.DisplayName))
	}

	t.SetSelectedFunc(func(row, col int) {
		j := uids[row]
		fn(j)
	})

	return t, nil
}

func newEntryTable(c *api.Client, key []byte, j *api.Journal) (*tview.Table, error) {
	t := tview.NewTable().SetSelectable(true, false)
	t.SetTitle("Entries").SetBorder(true)
	return t, nil
}

func Start(c *api.Client, key []byte) error {
	entries, err := newEntryTable(c, key, nil)
	if err != nil {
		return err
	}

	fn := func(j *api.Journal) {
		es, err := c.JournalEntries(j.UID)
		if err != nil {
			log.Fatal(err)
		}

		entries.Clear()
		for i, e := range es {
			content, err := e.GetContent(j, key)
			if err != nil {
				log.Fatal(err)
			}

			node, err := ical.ParseCalendar(content.Content)
			if err != nil {
				log.Fatal(err)
			}

			switch node.Name {
			case "VCARD":
				entries.SetCellSimple(i, 0, node.PropString("FN", "<N/A>"))
				entries.SetCellSimple(i, 1, node.PropString("TEL", ""))
			case "VCALENDAR", "VTODO":
				child := node.ChildByName("VTODO")
				if child == nil {
					child = node.ChildByName("VEVENT")
				}

				if child != nil {
					entries.SetCellSimple(i, 0, child.PropString("SUMMARY", "X"))
					when := child.PropDate("DTSTAMP", time.Time{})
					entries.SetCellSimple(i, 1, when.String())
				}
			}
		}
	}
	journals, err := newJournalTable(c, key, fn)
	if err != nil {
		return err
	}

	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(journals, 0, 1, true).
		AddItem(entries, 0, 2, false)

	return app.SetRoot(flex, true).Run()
}
