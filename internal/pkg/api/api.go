package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// Server exposes RESTful API endpoints.
type Server struct {
	Router *httprouter.Router
	logger *log.Logger
}

// NewServer instantiates an HTTP server and builds a route table.
func NewServer() *Server {
	s := &Server{
		Router: httprouter.New(),
		logger: log.New(os.Stderr, "API: ", 0),
	}
	s.buildRouteTable()
	return s
}

func (s *Server) buildRouteTable() {
	s.Router.GET("/api/probe", s.handleProbe())
}

// encodeJSON takes content of any type (v) and encodes to the writer (w) in
// JSON format.
func (s *Server) encodeJSON(w http.ResponseWriter, v interface{}, l *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
	}
}

// handleProbe acts as the Kubernetes liveness & readiness probes.
//
// In most scenarios, liveness & readiness probes should probably do different
// things. For example, a readiness probe might include making a call to a
// database to ensure that functionality works. However for this exercise
// simply being able to get a 200 from an endpoint should suffice for both
// probes.
func (s *Server) handleProbe() httprouter.Handle {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		payload := response{Message: "Application healthy!"}
		s.encodeJSON(w, payload, s.logger)
	}
}
