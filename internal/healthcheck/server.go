package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Server exposes the health status over HTTP.
type Server struct {
	checker *Checker
	addr    string
	server  *http.Server
}

// NewServer creates a health-check HTTP server bound to addr (e.g. ":9090").
func NewServer(addr string, checker *Checker) *Server {
	s := &Server{checker: checker, addr: addr}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	s.server = &http.Server{Addr: addr, Handler: mux}
	return s
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown() error {
	return s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	status := s.checker.Status()
	w.Header().Set("Content-Type", "application/json")
	if !status.Healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
	}
}
