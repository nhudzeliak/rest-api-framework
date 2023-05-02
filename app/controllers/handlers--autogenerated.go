package controllers

import (
    "net/http"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/api"
)

var (
	// These are controller instances.
	postsController	= PostsController{}
)

// MustInitialize performs all the needed setup for controllers.
func MustInitialize() { 
	postsController.MustInitialize()
}

var (
	// Controllers is a map of routes and functions that control them.
	Controllers = map[string]map[string]api.Serve { 
		"/posts": { 
			"GET": func(writer http.ResponseWriter, request *http.Request) {
		    	postsController.NewRequest(writer, request)
				postsController.IndexPosts()
			},
			"POST": func(writer http.ResponseWriter, request *http.Request) {
		    	postsController.NewRequest(writer, request)
				postsController.CreatePost()
			},
		},
		"/posts/{id}": { 
			"DELETE": func(writer http.ResponseWriter, request *http.Request) {
		    	postsController.NewRequest(writer, request)
				postsController.DeletePost()
			},
			"GET": func(writer http.ResponseWriter, request *http.Request) {
		    	postsController.NewRequest(writer, request)
				postsController.FindPost()
			},
			"PUT": func(writer http.ResponseWriter, request *http.Request) {
		    	postsController.NewRequest(writer, request)
				postsController.UpdatePost()
			},
		},
	}
)