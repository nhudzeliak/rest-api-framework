package api

import (
	"net/http"
)

// Serve is an alias to a http endpoint functional handler.
type Serve func(http.ResponseWriter, *http.Request)
