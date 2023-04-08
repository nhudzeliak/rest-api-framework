package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// APIClient implements some quality of life methods used by integration tests.
type APIClient struct {
	basePath string
}

// NewAPIClient instantiates APIClient with provided base path.
func NewAPIClient(path string) *APIClient {
	return &APIClient{basePath: path}
}

// Get ...
func (c *APIClient) Get(route string) (*http.Response, error) {
	return http.Get(c.basePath + route)
}

// PutBytes  ...
func (c *APIClient) PutBytes(route string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPut, c.basePath+route, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", http.DetectContentType(body))
	return http.DefaultClient.Do(request)
}

// PutObject  ...
func (c *APIClient) PutObject(route string, object interface{}) (*http.Response, error) {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return c.PutBytes(route, marshalled)
}

// PostBytes ...
func (c *APIClient) PostBytes(route string, body []byte) (*http.Response, error) {
	return http.Post(c.basePath+route, http.DetectContentType(body), bytes.NewBuffer(body))
}

// PostObject ...
func (c *APIClient) PostObject(route string, object interface{}) (*http.Response, error) {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return c.PostBytes(route, marshalled)
}

// Delete ...
func (c *APIClient) Delete(route string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodDelete, c.basePath+route, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(request)
}

// ParseJSONBody ...
func ParseJSONBody(body io.Reader, target any) error {
	return json.NewDecoder(body).Decode(&target)
}
