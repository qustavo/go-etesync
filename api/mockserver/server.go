package testserver

import (
	"net/http"
	"net/http/httptest"

	"github.com/gchaincl/go-etesync/api"
)

type Server struct {
	DB map[*api.Journal]api.Entries
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) Listen() func() {
	srv := httptest.NewServer(s)
	return srv.Close
}
