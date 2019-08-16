package main

import (
	"errors"
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gchaincl/go-etesync/api"
	"github.com/gchaincl/go-etesync/cache"
	"github.com/gchaincl/go-etesync/crypto"
	"github.com/gchaincl/go-etesync/gui"
	"github.com/gchaincl/go-etesync/store/sql"
	"github.com/laurent22/ical-go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

// Conf are the global flags
type Conf struct {
	url      string
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
		Version: "0.0.5",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "url", Usage: "Server URL", EnvVar: "ETESYNC_URL", Value: api.APIUrl, Destination: &cfg.url},
			cli.StringFlag{Name: "email", Usage: "login email", EnvVar: "ETESYNC_EMAIL", Destination: &cfg.email},
			cli.StringFlag{Name: "password", Usage: "login password", EnvVar: "ETESYNC_PASSWORD", Destination: &cfg.password},
			cli.StringFlag{Name: "key", Usage: "encryption key", EnvVar: "ETESYNC_KEY", Destination: &cfg.key},
			cli.StringFlag{Name: "db", Usage: "DB file path", Value: "~/.etecli.db", EnvVar: "ETESYNC_DB", Destination: &cfg.db},
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
					c, err := newCacheFromCtx(ctx)
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
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() != 1 {
						return errors.New("missing [uid]")
					}

					c, err := newCacheFromCtx(ctx)
					if err != nil {
						return nil
					}

					uid := ctx.Args()[0]
					return ete.JournalEntries(c, uid)
				},
			},
			cli.Command{
				Name: "gui", Usage: "Interactive gui",
				Action: func(ctx *cli.Context) error {
					cache, err := newCacheFromCtx(ctx)
					if err != nil {
						return err
					}
					return ete.StartGUI(cache)
				},
			},
		},
	}

	ete.runFn = app.RunAndExitOnError

	return ete
}

func newCacheFromCtx(ctx *cli.Context) (*cache.Cache, error) {
	client, err := newClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	store, err := newSQLStoreFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if err := store.Migrate(); err != nil {
		return nil, err
	}

	c := cache.New(store, client)
	if err := c.Sync(); err != nil {
		return nil, err
	}

	return c, nil

}

func newClientFromCtx(ctx *cli.Context) (*api.HTTPClient, error) {
	email := ctx.GlobalString("email")
	cl, err := api.NewClientWithURL(email, ctx.GlobalString("password"), ctx.GlobalString("url"))
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func newSQLStoreFromCtx(ctx *cli.Context) (*sql.Store, error) {
	db := expandPath(ctx.GlobalString("db"))
	store, err := sql.NewStore("sqlite3", db)
	if err != nil {
		return nil, err
	}

	if err := store.Migrate(); err != nil {
		return nil, err
	}

	return store, nil
}

func expandPath(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		return dir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(dir, path[2:])
	}
	return path
}

func (ete *EteCli) Journals(c *cache.Cache) error {
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

func (ete *EteCli) JournalEntries(c *cache.Cache, uid string) error {
	es, err := c.JournalEntries(uid)
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

		fmt.Printf("VCard %s\n", node)
	}

	return nil
}

func (ete *EteCli) StartGUI(cache *cache.Cache) error {
	return gui.New(cache, ete.key).Start()
}

func (ete *EteCli) Run() { ete.runFn() }

func main() {
	New().Run()
}
