package controllers

import (
	"context"
	"strconv"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/api"
	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/database"
	eposts "github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/entities"
	iposts "github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/interfaces"
	posts "github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/logic"
)

var ()

// PostsController is a wrapper for controllers that interact with posts.
type PostsController struct {
	api.ControllerSuite
	service iposts.PostsService
}

// MustInitialize performs all the setup needed for the controller.
func (c *PostsController) MustInitialize() {
	ctx := context.Background()
	reader, err := database.GetReader(ctx)
	if err != nil {
		panic(err)
	}
	writer, err := database.GetWriter(ctx)
	if err != nil {
		panic(err)
	}
	service, err := posts.NewPostsService(ctx, reader, writer)
	if err != nil {
		panic(err)
	}
	c.service = service
}

// IndexPosts fetches all posts.
func (c *PostsController) IndexPosts() {
	ctx := context.Background()
	ps, err := c.service.IndexPosts(ctx)
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	c.ServeOK(ps)
}

// FindPost fetches a single post.
func (c *PostsController) FindPost() {
	ctx := context.Background()
	rawID := c.ParseURLParams()["id"]
	id, err := strconv.Atoi(rawID)
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	p, err := c.service.FindPost(ctx, eposts.PostID(id))
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	c.ServeOK(p)
}

// UpdatePost updates a post.
func (c *PostsController) UpdatePost() {
	ctx := context.Background()
	rawID := c.ParseURLParams()["id"]
	id, err := strconv.Atoi(rawID)
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	var p eposts.Post
	err = c.ParseJSONBody(&p)
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	p.ID = eposts.PostID(id)
	err = c.service.UpdatePost(ctx, &p)
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	c.ServeCreated(p)
}

// CreatePost creates a post.
func (c *PostsController) CreatePost() {
	ctx := context.Background()
	var p eposts.Post
	err := c.ParseJSONBody(&p)
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	err = c.service.CreatePost(ctx, &p)
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	c.ServeCreated(p)
}

// DeletePost deletes a post.
func (c *PostsController) DeletePost() {
	ctx := context.Background()
	rawID := c.ParseURLParams()["id"]
	id, err := strconv.Atoi(rawID)
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	err = c.service.DeletePost(ctx, eposts.PostID(id))
	if err != nil {
		c.ServeBadRequest(err.Error())
		return
	}
	c.ServeMessageOK("post deleted")
}
