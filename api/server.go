package api

import (
	"log"
	"net/http"
	"sync"

	"github.com/ebastien/mznapi/solver"
	"github.com/google/uuid"
)

// Server maintains the state of the HTTP APIs.
type Server struct {
	models  map[uuid.UUID]solver.Model
	address string
	baseURL string
	router  http.Handler
	lock    sync.RWMutex
	workers chan struct{}
}

// NewServer creates a new server instance.
func NewServer(addr string, parallelism int) *Server {
	server := &Server{
		models:  make(map[uuid.UUID]solver.Model),
		address: addr,
		baseURL: "http://" + addr,
		workers: make(chan struct{}, parallelism),
	}
	server.router = newRouter(server)
	return server
}

// Serve runs the server main loop to handle incoming connections.
func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(s.address, s.router))
}

// ServeHTTP proxies to the underlying router implementation.
// It implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
