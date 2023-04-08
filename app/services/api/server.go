package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server wraps http.Server.
type Server struct {
	host        string
	port        int
	multiplexer http.Handler
}

// NewServer instantiates a new Server.
func NewServer(host string, port int, controllers map[string]map[string]Serve) *Server {
	multiplexer := mux.NewRouter()
	for route, methodGroup := range controllers {
		serve := generalizeHandler(methodGroup)
		multiplexer.Handle(route, handler{serve: serve})
	}
	return &Server{
		host:        host,
		port:        port,
		multiplexer: multiplexer,
	}
}

func (s *Server) Run() error {
	return http.ListenAndServe(fmt.Sprintf("%v:%v", s.host, s.port), s.multiplexer)
}

// generalizeHandler wraps Serve for all provided methods into one function.
func generalizeHandler(handlers map[string]Serve) Serve {
	return func(writer http.ResponseWriter, request *http.Request) {
		if handler, ok := handlers[request.Method]; ok {
			handler(writer, request)
		} else {
			suite := &ControllerSuite{
				writer:  writer,
				request: request,
			}
			suite.ServeNotFound()
		}
	}
}
