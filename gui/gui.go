package gui

import (
	"log"
	"time"

	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/cache"
	"github.com/gchaincl/go-etesync/crypto"
	"github.com/gdamore/tcell"
	"github.com/kofoworola/godate"
	"github.com/laurent22/ical-go"
	"github.com/rivo/tview"
)

type GUI struct {
	app      *tview.Application
	page     *tview.Pages
	entries  *tview.Table
	journals *tview.Table

	cache *cache.Cache
	key   []byte
}

func New(cache *cache.Cache, key []byte) *GUI {
	gui := &GUI{
		app:   tview.NewApplication(),
		cache: cache,
		key:   key,
	}

	gui.page = tview.NewPages()
	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			gui.page.HidePage("help")
			gui.page.ShowPage("flex")
		})
	gui.page.AddPage("help", modal, true, false)

	gui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'h', '?':
				gui.page.HidePage("flex")
				gui.page.ShowPage("help")
			}
		}
		return event
	})

	return gui
}

func (gui *GUI) newEntries() *tview.Table {
	t := tview.NewTable().SetSelectable(true, false).SetFixed(1, 1)
	t.SetBorder(true)

	t.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Key() {
		case tcell.KeyLeft, tcell.KeyTAB:
			gui.app.SetFocus(gui.journals)
		}

		return e
	})

	return t
}

func (gui *GUI) newJournals() (*tview.Table, error) {
	js, err := gui.cache.Journals()
	if err != nil {
		return nil, err
	}

	t := tview.NewTable().SetSelectable(true, false)
	t.SetTitle("Journals").SetBorder(true)

	uids := make([]*api.Journal, len(js))
	for i, j := range js {
		cipher := crypto.New([]byte(j.UID), gui.key)
		content, err := j.GetContent(cipher)
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
		err := gui.onJournalSelect(j)
		if err != nil {
			log.Fatal(err)
		}
	})

	return t, nil
}

func setTableHeaders(t *tview.Table, headers ...string) {
	for i, s := range headers {
		cell := tview.NewTableCell(s).
			SetSelectable(false).
			SetTextColor(tcell.ColorGray)
		t.SetCell(0, i, cell)
	}
}

func (gui *GUI) onJournalSelect(j *api.Journal) error {
	es, err := gui.cache.JournalEntries(j.UID)
	if err != nil {
		return err
	}

	cipher := crypto.New([]byte(j.UID), gui.key)
	jc, err := j.GetContent(cipher)
	if err != nil {
		log.Fatal(err)
	}
	gui.entries.SetTitle(string(jc.Type))
	gui.app.SetFocus(gui.entries)

	gui.entries.Clear()
	for i := 0; i < len(es); i++ {
		// as entries are sorted from older to newer we get them from newer to older
		e := es[len(es)-i-1]

		content, err := e.GetContent(cipher)
		if err != nil {
			return err
		}

		node, err := ical.ParseCalendar(content.Content)
		if err != nil {
			return err
		}

		var icon string
		switch content.Action {
		case "ADD":
			icon = "✔"
		case "DELETE":
			icon = "✖"
		case "CHANGE":
			icon = "↪"
		default:
			icon = content.Action
		}
		switch node.Name {
		case "VCARD":
			// set headers
			if i == 0 {
				setTableHeaders(gui.entries, "", "Name", "Phone")
			}

			gui.entries.SetCellSimple(i+1, 0, icon)
			gui.entries.SetCellSimple(i+1, 1, node.PropString("FN", "<N/A>"))
			gui.entries.SetCellSimple(i+1, 2, node.PropString("TEL", ""))
		case "VCALENDAR", "VTODO":
			// set headers
			if i == 0 {
				setTableHeaders(gui.entries, "", "Summary", "Date")
			}

			child := node.ChildByName("VTODO")
			if child == nil {
				child = node.ChildByName("VEVENT")
			}

			if child != nil {
				gui.entries.SetCellSimple(i+1, 0, icon)
				gui.entries.SetCellSimple(i+1, 1, child.PropString("SUMMARY", ""))
				when := child.PropDate("DTSTAMP", time.Time{})
				diff := godate.Create(when).DifferenceFromNowForHumans()
				gui.entries.SetCellSimple(i+1, 2, diff)
			}
		default:
			panic(node.Name)
		}

		gui.entries.Select(1, 0)
		gui.entries.ScrollToBeginning()
	}
	return nil
}

func (gui *GUI) Start() error {
	gui.entries = gui.newEntries()

	var err error
	gui.journals, err = gui.newJournals()
	if err != nil {
		return err
	}

	flex := tview.NewFlex().
		AddItem(gui.journals, 0, 1, true).
		AddItem(gui.entries, 0, 2, false)

	gui.page.AddAndSwitchToPage("flex", flex, true)
	return gui.app.SetRoot(gui.page, true).Run()
}
