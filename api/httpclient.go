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

var _ Client = &HTTPClient{}

// HTTPClient is a EteSync API Client
type HTTPClient struct {
	apiurl   string
	username string
	password string
	token    string
	debug    bool
}

// NewClient returns a new HTTPClient given a username and password
func NewClient(u, p string) (*HTTPClient, error) {
	c := &HTTPClient{
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
func (c *HTTPClient) WithDebug() *HTTPClient {
	n := *c
	n.debug = true
	return &n
}

func (c *HTTPClient) url(paths ...string) string {
	s := append([]string{c.apiurl}, paths...)
	return strings.Join(s, "/")
}

func (c *HTTPClient) withHeaders(req *http.Request) *http.Request {
	req.Header.Set("Content-Type", "application/json")
	if t := c.token; t != "" {
		req.Header.Set("Authorization", "Token "+t)
	}
	return req
}

func (c *HTTPClient) post(path string, src, dst interface{}) (int, error) {
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

func (c *HTTPClient) auth() error {
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

func (c *HTTPClient) get(path string, dst interface{}) (int, error) {
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
func (c *HTTPClient) Journals() (Journals, error) {
	dst := Journals{}

	if _, err := c.get("api/v1/journals/", &dst); err != nil {
		return nil, err
	}

	return dst, nil
}

func (c *HTTPClient) Journal(uid string) (*Journal, error) {
	dst := Journal{}
	if _, err := c.get("api/v1/journals/"+uid, &dst); err != nil {
		return nil, err
	}

	return &dst, nil
}

func (c *HTTPClient) JournalEntries(uid string) (Entries, error) {
	dst := Entries{}
	if _, err := c.get("api/v1/journals/"+uid+"/entries", &dst); err != nil {
		return nil, err
	}

	return dst, nil
}
