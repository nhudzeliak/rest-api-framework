package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// handler is a wrapper for Serve that implements http.Handler and is used purely for adaptation purposes.
type handler struct {
	serve Serve
}

// ServeHTTP wraps Serve.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("%-7v\t%v", r.Method, r.URL)
	h.serve(w, r)
}
