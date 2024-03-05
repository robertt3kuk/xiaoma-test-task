// Package httpserver implements HTTP server.
package httpserver

import (
	"time"

	"github.com/valyala/fasthttp"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server          *fasthttp.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New -.
func New(handler fasthttp.RequestHandler, port string, opts ...Option) *Server {
	httpServer := &fasthttp.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start(port)

	return s
}

func (s *Server) start(port string) {
	if port == "" {
		port = _defaultAddr
	}
	go func() {
		s.notify <- s.server.ListenAndServe(port)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {

	return s.server.Shutdown()
}
