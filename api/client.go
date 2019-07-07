package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

var (
	// ErrInvalidCredentials denotes invalid credentials when trying to get a API token
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// APIUrl is the default URL
const APIUrl = "https://api.etesync.com"

// Client is a EteSync API Client
type Client struct {
	apiurl   string
	username string
	password string
	token    string
	debug    bool
}

// NewClient returns a new client given a username and password
func NewClient(u, p string) (*Client, error) {
	c := &Client{
		apiurl:   APIUrl,
		username: u,
		password: p,
	}

	if err := c.auth(); err != nil {
		return nil, err
	}

	return c, nil
}

// WithDebug returns a copy of the client with debug enabled
// When debug is enabled HTTP request and responses are logged
func (c *Client) WithDebug() *Client {
	n := *c
	n.debug = true
	return &n
}

func (c *Client) url(paths ...string) string {
	s := append([]string{c.apiurl}, paths...)
	return strings.Join(s, "/")
}

func (c *Client) withHeaders(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	if t := c.token; t != "" {
		req.Header.Set("Authorization", "Token "+t)
	}
	return req
}

func (c *Client) post(path string, src, dst interface{}) (int, error) {
	buf, err := json.Marshal(src)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", c.url(path), bytes.NewBuffer(buf))
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(c.withHeaders(req))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func (c *Client) auth() error {
	src := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{c.username, c.password}

	dst := struct {
		Token string `json:"token"`
	}{}

	status, err := c.post("api-token-auth/", src, &dst)
	if err != nil {
		return err
	}

	if status == 400 {
		return ErrInvalidCredentials
	}

	c.token = dst.Token

	return nil
}

func (c *Client) get(path string, dst interface{}) (int, error) {
	req, err := http.NewRequest("GET", c.url(path), nil)
	if err != nil {
		return 0, err
	}

	if c.debug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return 0, err
		}
		fmt.Println(string(dump))
	}

	resp, err := http.DefaultClient.Do(c.withHeaders(req))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if c.debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return 0, err
		}
		fmt.Println(string(dump))
	}

	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

// Journals retrieves the available journals
func (c *Client) Journals() (Journals, error) {
	dst := Journals{}

	if _, err := c.get("api/v1/journals/", &dst); err != nil {
		return nil, err
	}

	return dst, nil
}

func (c *Client) Journal(uid string) (Entries, error) {
	dst := Entries{}
	if _, err := c.get("api/v1/journals/"+uid+"/entries", &dst); err != nil {
		return nil, err
	}

	return dst, nil
}
