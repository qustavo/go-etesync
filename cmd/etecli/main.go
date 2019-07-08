package main

import (
	"errors"
	"fmt"

	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/gui"
	"github.com/laurent22/ical-go"
	"github.com/urfave/cli"
)

type App struct {
	cli    *cli.App
	client *api.Client
}

func NewApp() *App {
	app := &App{}

	app.cli = cli.NewApp()
	app.cli.Version = "0.0.1"
	app.cli.Name = "etecli"
	app.cli.Usage = "ETESync cli tool"
	app.cli.Flags = []cli.Flag{
		cli.StringFlag{Name: "email", Usage: "login email", EnvVar: "ETESYNC_EMAIL"},
		cli.StringFlag{Name: "password", Usage: "login password", EnvVar: "ETESYNC_EMAIL"},
		cli.StringFlag{Name: "key", Usage: "Encryption key", EnvVar: "ETESYNC_KEY"},
	}

	app.cli.Commands = []cli.Command{
		cli.Command{
			Name: "journals", Usage: "Display available journals", Category: "api",
			Action: func(ctx *cli.Context) error {
				c, err := newClientFromCtx(ctx)
				if err != nil {
					return nil
				}
				key := []byte(ctx.GlobalString("key"))
				return Journals(c, key)
			},
		},
		cli.Command{
			Name: "journal", Usage: "Retrieve a journal given a uid", Category: "api", ArgsUsage: "[uid]",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("missing [uid]")
				}

				c, err := newClientFromCtx(ctx)
				if err != nil {
					return nil
				}

				uid := ctx.Args()[0]
				key := []byte(ctx.GlobalString("key"))

				return Journal(c, uid, key)
			},
		},

		cli.Command{
			Name: "entries", Usage: "displays entries given a journal uid", Category: "api", ArgsUsage: "[uid]",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("missing [uid]")
				}

				c, err := newClientFromCtx(ctx)
				if err != nil {
					return nil
				}

				uid := ctx.Args()[0]
				key := []byte(ctx.GlobalString("key"))

				return JournalEntries(c, uid, key)
			},
		},
		cli.Command{
			Name: "gui", Usage: "Interactive gui",
			Action: func(ctx *cli.Context) error {
				c, err := newClientFromCtx(ctx)
				if err != nil {
					return err
				}

				key := []byte(ctx.GlobalString("key"))
				return StartGUI(c, key)
			},
		},
	}

	return app
}

func newClientFromCtx(ctx *cli.Context) (*api.Client, error) {
	return api.NewClient(ctx.GlobalString("email"), ctx.GlobalString("password"))
}

func Journals(c *api.Client, key []byte) error {
	js, err := c.Journals()
	if err != nil {
		return err
	}

	for _, j := range js {
		content, err := j.GetContent(key)
		if err != nil {
			return err
		}
		fmt.Printf("<Journal uid:%s\n     content: %s>\n", j.UID, content)
	}

	return nil
}

func Journal(c *api.Client, uid string, key []byte) error {
	j, err := c.Journal(uid)
	if err != nil {
		return err
	}

	content, err := j.GetContent(key)
	if err != nil {
		return err
	}

	fmt.Printf("content  :%s\n", content)
	fmt.Printf("owner    :%s\n", j.Owner)
	fmt.Printf("read-only:%v\n", j.ReadOnly)

	return nil
}

func JournalEntries(c *api.Client, uid string, key []byte) error {
	j, err := c.Journal(uid)
	if err != nil {
		return err
	}

	es, err := c.JournalEntries(j.UID)
	if err != nil {
		return err
	}

	for _, e := range es {
		content, err := e.GetContent(j, key)
		if err != nil {
			return err
		}

		fmt.Printf("UID: %s\n", e.UID)
		node, err := ical.ParseCalendar(content.Content)
		if err != nil {
			return err

		}
		c := node.ChildByName("VEVENT")
		fmt.Printf("%s", c.PropString("SUMMARY", "XXX"))
		fmt.Printf("-----\n")
	}

	return nil
}

func StartGUI(c *api.Client, key []byte) error {
	return gui.Start(c, key)
}

func (app *App) Run() { app.cli.RunAndExitOnError() }

func main() {
	NewApp().Run()
}
