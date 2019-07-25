package main

import (
	"errors"
	"fmt"

	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/crypto"
	"github.com/gchaincl/go-etesync/gui"
	"github.com/gchaincl/go-etesync/store"
	"github.com/gchaincl/go-etesync/store/sql"
	"github.com/laurent22/ical-go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

// Conf are the global flags
type Conf struct {
	server   string
	email    string
	password string
	key      string
	db       string
	sync     bool
}

type EteCli struct {
	cfg   *Conf
	key   []byte
	runFn func()
}

func New() *EteCli {
	cfg := &Conf{}
	ete := &EteCli{cfg: cfg}

	app := &cli.App{
		Name:    "etecli",
		Usage:   "ETESync cli tool",
		Version: "0.0.1",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "server", Usage: "Server URL", EnvVar: "ETESYNC_SERVER", Value: "", Destination: &cfg.server},
			cli.StringFlag{Name: "email", Usage: "login email", EnvVar: "ETESYNC_EMAIL", Destination: &cfg.email},
			cli.StringFlag{Name: "password", Usage: "login password", EnvVar: "ETESYNC_PASSWORD", Destination: &cfg.password},
			cli.StringFlag{Name: "key", Usage: "encryption key", EnvVar: "ETESYNC_KEY", Destination: &cfg.key},
			cli.StringFlag{Name: "db", Usage: "DB file path", Value: "~/.etecli.db", EnvVar: "ETESYNC_DB", Destination: &cfg.db},
			cli.BoolFlag{Name: "sync", Usage: "force sync on start", Destination: &cfg.sync},
		},

		Before: func(ctx *cli.Context) error {
			if cfg.email == "" {
				return errors.New("missing `--email` flag")
			}

			if cfg.password == "" {
				return errors.New("missing `--password` flag")
			}

			if cfg.key == "" {
				return errors.New("missing `--key` flag")
			}

			var err error
			ete.key, err = api.DeriveKey(cfg.email, []byte(cfg.key))
			if err != nil {
				return err
			}
			return nil
		},

		Commands: []cli.Command{
			cli.Command{
				Name: "journals", Usage: "Display available journals", Category: "api",
				Action: func(ctx *cli.Context) error {
					c, err := newClientFromCtx(ctx)
					if err != nil {
						return nil
					}
					return ete.Journals(c)
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
						return err
					}

					uid := ctx.Args()[0]
					return ete.Journal(c, uid)
				},
			},

			cli.Command{
				Name: "entries", Usage: "displays entries given a journal uid", Category: "api", ArgsUsage: "[uid]",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "last", Usage: "get entries after <last> uid"},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() != 1 {
						return errors.New("missing [uid]")
					}

					c, err := newClientFromCtx(ctx)
					if err != nil {
						return nil
					}

					uid := ctx.Args()[0]
					last := ctx.String("last")
					return ete.JournalEntries(c, uid, last)
				},
			},
			cli.Command{
				Name: "gui", Usage: "Interactive gui",
				Action: func(ctx *cli.Context) error {
					c, err := newClientFromCtx(ctx)
					if err != nil {
						return err
					}

					s, err := sql.NewStore("sqlite3", "/tmp/etesync.db")
					if err != nil {
						return err
					}

					if err := s.Migrate(); err != nil {
						return err
					}

					return ete.StartGUI(c, s)
				},
			},
		},
	}

	ete.runFn = app.RunAndExitOnError

	return ete
}

func newClientFromCtx(ctx *cli.Context) (*api.HTTPClient, error) {
	email := ctx.GlobalString("email")
	cl, err := api.NewClientWithURL(email, ctx.GlobalString("password"), ctx.GlobalString("server"))
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (ete *EteCli) Journals(c api.Client) error {
	js, err := c.Journals()
	if err != nil {
		return err
	}

	for _, j := range js {
		fmt.Printf("<Journal uid:%s>\n", j.UID)
	}

	return nil
}

func (ete *EteCli) Journal(c api.Client, uid string) error {
	j, err := c.Journal(uid)
	if err != nil {
		return err
	}

	cipher := crypto.New([]byte(uid), ete.key)
	content, err := j.GetContent(cipher)
	if err != nil {
		return err
	}

	fmt.Printf("name     : %s\n", content.DisplayName)
	fmt.Printf("type     : %s\n", content.Type)
	fmt.Printf("owner    : %s\n", j.Owner)
	fmt.Printf("read-only: %v\n", j.ReadOnly)

	return nil
}

func (ete *EteCli) JournalEntries(c api.Client, uid string, last string) error {
	var arg *string = nil
	if last != "" {
		arg = &last
	}

	es, err := c.JournalEntries(uid, arg)
	if err != nil {
		return err
	}

	cipher := crypto.New([]byte(uid), ete.key)

	for _, e := range es {
		content, err := e.GetContent(cipher)
		if err != nil {
			return err
		}

		fmt.Printf("UID: %s\n", e.UID)
		node, err := ical.ParseCalendar(content.Content)
		if err != nil {
			return err

		}

		fmt.Printf("VCard %s", node)
	}

	return nil
}

func (ete *EteCli) StartGUI(c api.Client, s store.Store) error {
	return gui.Start(c, s, ete.key)
}

func (ete *EteCli) Run() { ete.runFn() }

func main() {
	New().Run()
}
