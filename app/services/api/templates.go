package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// ControllerSuite implements some quality-of-life functionality.
type ControllerSuite struct {
	writer  http.ResponseWriter
	request *http.Request
}

// NewRequest resets the controller.
func (s *ControllerSuite) NewRequest(writer http.ResponseWriter, request *http.Request) {
	s.writer = writer
	s.request = request
}

// ServeOK serves a 200 response.
func (s *ControllerSuite) ServeOK(object interface{}) {
	s.writer.WriteHeader(http.StatusOK)
	s.RenderJSON(object)
}

// ServeMessageOK serves a 200 response with string message.
func (s *ControllerSuite) ServeMessageOK(message string) {
	response := make(map[string]string)
	response["message"] = message
	s.ServeOK(response)
}

// ServeEmptyOK serves an empty 200 response.
func (s *ControllerSuite) ServeEmptyOK() {
	s.writer.WriteHeader(http.StatusOK)
}

// ServeCreated serves a 201 response with a provided object.
func (s *ControllerSuite) ServeCreated(object interface{}) {
	s.writer.WriteHeader(http.StatusCreated)
	s.RenderJSON(object)
}

// ServeBadRequest serves a 400 response with a provided message.
func (s *ControllerSuite) ServeBadRequest(message string) {
	s.writer.WriteHeader(http.StatusBadRequest)
	response := make(map[string]string)
	response["message"] = message
	s.RenderJSON(response)
}

// ServeNotFound serves a standard 404 error.
func (s *ControllerSuite) ServeNotFound() {
	s.writer.WriteHeader(http.StatusNotFound)
	response := make(map[string]string)
	response["message"] = fmt.Sprintf("no handler registered at route %v for method %v", s.request.Method, s.request.URL.Path)
	s.RenderJSON(response)
}

// ServeConflict serves a 409 response with a provided message.
func (s *ControllerSuite) ServeConflict(message string) {
	s.writer.WriteHeader(http.StatusConflict)
	response := make(map[string]string)
	response["message"] = message
	s.RenderJSON(response)
}

func (s *ControllerSuite) ServeInternalError(message string) {
	s.writer.WriteHeader(http.StatusInternalServerError)
	response := make(map[string]string)
	response["message"] = message
	s.RenderJSON(response)
}

// RenderJSON writes a json to response.
func (s *ControllerSuite) RenderJSON(response interface{}) {
	s.writer.Header().Set("Content-Type", "application/json")
	bytesResponse, err := json.Marshal(response)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to marshal response")
	}
	_, err = s.writer.Write(bytesResponse)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to write response")
	}
}

// ParseURLParams searches the url path for patterns and returns a dict of string/string pairs.
func (s *ControllerSuite) ParseURLParams() map[string]string {
	return mux.Vars(s.request)
}

// ParseJSONBody parses request body.
// This function really shouldn't be here...
func (s *ControllerSuite) ParseJSONBody(target any) error {
	return json.NewDecoder(s.request.Body).Decode(&target)
}
