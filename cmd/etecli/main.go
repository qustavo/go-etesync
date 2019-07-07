package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gchaincl/go-etesync/api"
	"github.com/urfave/cli"
)

type App struct {
	cli    *cli.App
	client *api.Client
}

func NewApp() *App {
	app := &App{}

	app.cli = cli.NewApp()
	app.cli.Name = "etecli"
	app.cli.Usage = "ETESync cli tool"
	app.cli.Flags = []cli.Flag{
		cli.StringFlag{Name: "email"},
		cli.StringFlag{Name: "password"},
		cli.StringFlag{Name: "key", Usage: "Encryption key"},
	}

	app.cli.Commands = []cli.Command{
		JournalsCmd(),
		JournalCmd(),
		EntriesCmd(),
	}

	return app
}

func newClientFromCtx(ctx *cli.Context) (*api.Client, error) {
	return api.NewClient(ctx.GlobalString("email"), ctx.GlobalString("password"))
}

func JournalsCmd() cli.Command {
	return cli.Command{
		Name:     "journals",
		Usage:    "Display available journals",
		Category: "api",
		Action: func(ctx *cli.Context) error {
			c, err := newClientFromCtx(ctx)
			if err != nil {
				return nil
			}

			js, err := c.Journals()
			if err != nil {
				return err
			}

			for _, j := range js {
				content, err := j.GetContent([]byte(ctx.GlobalString("key")))
				if err != nil {
					return err
				}
				fmt.Printf("<Journal uid:%s\n     content: %s>\n", j.UID, content)
			}

			return nil
		},
	}
}

func JournalCmd() cli.Command {
	return cli.Command{
		Name:      "journal",
		Usage:     "Retrieve a journal given a uid",
		Category:  "api",
		ArgsUsage: "[uid]",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("missing [uid]")
			}

			c, err := newClientFromCtx(ctx)
			if err != nil {
				return err
			}

			j, err := c.Journal(ctx.Args()[0])
			if err != nil {
				return err
			}

			content, err := j.GetContent([]byte(ctx.GlobalString("key")))
			if err != nil {
				return err
			}
			fmt.Printf("content  :%s\n", content)
			fmt.Printf("owner    :%s\n", j.Owner)
			fmt.Printf("read-only:%v\n", j.ReadOnly)

			return nil
		},
	}
}

func EntriesCmd() cli.Command {
	return cli.Command{
		Name:      "entries",
		Usage:     "displays entries given a journal uid",
		Category:  "api",
		ArgsUsage: "[uid]",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("missing [uid]")
			}

			c, err := newClientFromCtx(ctx)
			if err != nil {
				return err
			}

			j, err := c.Journal(ctx.Args()[0])
			if err != nil {
				return err
			}

			es, err := c.JournalEntries(j.UID)
			if err != nil {
				return err
			}

			for _, e := range es {
				content, err := e.GetContent(j, []byte(ctx.GlobalString("key")))
				if err != nil {
					return err
				}

				log.Printf("----\ncontent = %s-----\n", content)
			}

			return nil
		},
	}
}

func (app *App) Run() { app.cli.RunAndExitOnError() }

func main() {
	NewApp().Run()
}
