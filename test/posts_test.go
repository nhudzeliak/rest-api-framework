package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/entities"

	"github.com/stretchr/testify/assert"
)

func TestPostsController_FindPost(t *testing.T) {
	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(entities.PostID) (*http.Response, error)
		assertion func(*http.Response, error)
		cleanup   func(entities.PostID) error
	}{
		// Existing id.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   "test-title",
					Content: "test-content",
				}
				response, err := apiClient.PostObject("/posts", post)
				if err != nil {
					return 0, err
				}
				err = ParseJSONBody(response.Body, &post)
				if err != nil {
					return 0, err
				}
				return post.ID, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				return apiClient.Get(fmt.Sprintf("/posts/%v", id))
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusOK)
					var post entities.Post
					if err := ParseJSONBody(response.Body, &post); assert.NoError(t, err) {
						assert.NotZero(t, post)
					}
				}
			},
			cleanup: func(id entities.PostID) error {
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				return err
			},
		},
		// Non-existing id.
		{
			setup: func() (entities.PostID, error) {
				id := entities.PostID(69)
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				if err != nil {
					return 0, err
				}
				return id, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				return apiClient.Get(fmt.Sprintf("/posts/%v", id))
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusBadRequest)
					var message map[string]string
					if err := ParseJSONBody(response.Body, &message); assert.NoError(t, err) {
						if m, ok := message["message"]; assert.True(t, ok) {
							assert.Equal(t, m, entities.ErrPostNotFound.Error())
						}
					}
				}
			},
			cleanup: func(id entities.PostID) error { return nil },
		},
	}

	for _, c := range cases {
		id, err := c.setup()
		assert.Nil(t, err)
		c.assertion(c.base(id))
		err = c.cleanup(id)
		assert.Nil(t, err)
	}
}

func TestPostsController_IndexPosts(t *testing.T) {
	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(entities.PostID) (*http.Response, error)
		assertion func(*http.Response, error)
		cleanup   func(entities.PostID) error
	}{
		// Non-empty table.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   "test-title",
					Content: "test-content",
				}
				response, err := apiClient.PostObject("/posts", post)
				if err != nil {
					return 0, err
				}
				err = ParseJSONBody(response.Body, &post)
				if err != nil {
					return 0, err
				}
				return post.ID, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				return apiClient.Get("/posts")
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusOK)
					var posts []entities.Post
					if err := ParseJSONBody(response.Body, &posts); assert.NoError(t, err) {
						if assert.NotEmpty(t, posts) {
							assert.NotZero(t, posts[0])
						}
					}
				}
			},
			cleanup: func(id entities.PostID) error {
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				return err
			},
		},
	}

	for _, c := range cases {
		id, err := c.setup()
		assert.Nil(t, err)
		c.assertion(c.base(id))
		err = c.cleanup(id)
		assert.Nil(t, err)
	}
}

func TestPostsController_UpdatePost(t *testing.T) {
	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(entities.PostID) (*http.Response, error)
		assertion func(*http.Response, error)
		cleanup   func(entities.PostID) error
	}{
		// Existing id.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   "test-title",
					Content: "test-content",
				}
				response, err := apiClient.PostObject("/posts", post)
				if err != nil {
					return 0, err
				}
				err = ParseJSONBody(response.Body, &post)
				if err != nil {
					return 0, err
				}
				return post.ID, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				post := entities.Post{
					Title:   "test-title-new",
					Content: "test-content-new",
				}
				return apiClient.PutObject(fmt.Sprintf("/posts/%v", id), post)
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusCreated)
					var post entities.Post
					if err := ParseJSONBody(response.Body, &post); assert.NoError(t, err) {
						assert.NotZero(t, post)
					}
				}
			},
			cleanup: func(id entities.PostID) error {
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				return err
			},
		},
		// Non-existing id.
		{
			setup: func() (entities.PostID, error) {
				id := entities.PostID(69)
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				if err != nil {
					return 0, err
				}
				return id, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				post := entities.Post{
					Title:   "test-title-new",
					Content: "test-content-new",
				}
				return apiClient.PutObject(fmt.Sprintf("/posts/%v", id), post)
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusBadRequest)
					var message map[string]string
					if err := ParseJSONBody(response.Body, &message); assert.NoError(t, err) {
						if m, ok := message["message"]; assert.True(t, ok) {
							assert.Equal(t, m, entities.ErrPostNotFound.Error())
						}
					}
				}
			},
			cleanup: func(id entities.PostID) error { return nil },
		},
		// Invalid post.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   "test-title",
					Content: "test-content",
				}
				response, err := apiClient.PostObject("/posts", post)
				if err != nil {
					return 0, err
				}
				err = ParseJSONBody(response.Body, &post)
				if err != nil {
					return 0, err
				}
				return post.ID, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 256),
					Content: "test-content-new",
				}
				return apiClient.PutObject(fmt.Sprintf("/posts/%v", id), post)
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusBadRequest)
					var message map[string]string
					if err := ParseJSONBody(response.Body, &message); assert.NoError(t, err) {
						if m, ok := message["message"]; assert.True(t, ok) {
							assert.Equal(t, m, entities.ErrInvalidTitle.Error())
						}
					}
				}
			},
			cleanup: func(id entities.PostID) error {
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				return err
			},
		},
	}

	for _, c := range cases {
		id, err := c.setup()
		assert.Nil(t, err)
		c.assertion(c.base(id))
		err = c.cleanup(id)
		assert.Nil(t, err)
	}
}

func TestPostsController_CreatePost(t *testing.T) {
	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(entities.PostID) (*http.Response, error)
		assertion func(*http.Response, error)
		cleanup   func(entities.PostID) error
	}{
		// Valid id.
		{
			setup: func() (entities.PostID, error) { return 0, nil },
			base: func(id entities.PostID) (*http.Response, error) {
				post := entities.Post{
					Title:   "test-title-new",
					Content: "test-content-new",
				}
				return apiClient.PostObject("/posts", post)
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusCreated)
					var post entities.Post
					if err := ParseJSONBody(response.Body, &post); assert.NoError(t, err) {
						assert.NotZero(t, post)
					}
				}
			},
			cleanup: func(id entities.PostID) error {
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				return err
			},
		},
		// Conflicting id.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   "test-title",
					Content: "test-content",
				}
				response, err := apiClient.PostObject("/posts", post)
				if err != nil {
					return 0, err
				}
				err = ParseJSONBody(response.Body, &post)
				if err != nil {
					return 0, err
				}
				return post.ID, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				post := entities.Post{
					ID:      id,
					Title:   "test-title-new",
					Content: "test-content-new",
				}
				return apiClient.PostObject("/posts", post)
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusBadRequest)
					var message map[string]string
					if err := ParseJSONBody(response.Body, &message); assert.NoError(t, err) {
						if m, ok := message["message"]; assert.True(t, ok) {
							assert.Equal(t, m, entities.ErrDuplicatePost.Error())
						}
					}
				}
			},
			cleanup: func(id entities.PostID) error {
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				return err
			},
		},
		// Invalid post.
		{
			setup: func() (entities.PostID, error) { return 0, nil },
			base: func(id entities.PostID) (*http.Response, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 256),
					Content: "test-content-new",
				}
				return apiClient.PostObject("/posts", post)
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusBadRequest)
					var message map[string]string
					if err := ParseJSONBody(response.Body, &message); assert.NoError(t, err) {
						if m, ok := message["message"]; assert.True(t, ok) {
							assert.Equal(t, m, entities.ErrInvalidTitle.Error())
						}
					}
				}
			},
			cleanup: func(id entities.PostID) error {
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				return err
			},
		},
	}

	for _, c := range cases {
		id, err := c.setup()
		assert.Nil(t, err)
		response, err := c.base(id)
		c.assertion(response, err)
		var post entities.Post
		if err := ParseJSONBody(response.Body, &post); err == nil {
			err = c.cleanup(post.ID)
			assert.Nil(t, err)
			continue
		}
		err = c.cleanup(id)
		assert.Nil(t, err)
	}
}

func TestPostsController_DeletePost(t *testing.T) {
	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(entities.PostID) (*http.Response, error)
		assertion func(*http.Response, error)
	}{
		// Existing id.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   "test-title",
					Content: "test-content",
				}
				response, err := apiClient.PostObject("/posts", post)
				if err != nil {
					return 0, err
				}
				err = ParseJSONBody(response.Body, &post)
				if err != nil {
					return 0, err
				}
				return post.ID, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				return apiClient.Delete(fmt.Sprintf("/posts/%v", id))
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusOK)
					var message map[string]string
					if err := ParseJSONBody(response.Body, &message); assert.NoError(t, err) {
						if m, ok := message["message"]; assert.True(t, ok) {
							assert.Equal(t, m, "post deleted")
						}
					}
				}
			},
		},
		// Non-existing id.
		{
			setup: func() (entities.PostID, error) {
				id := entities.PostID(69)
				_, err := apiClient.Delete(fmt.Sprintf("/posts/%v", id))
				if err != nil {
					return 0, err
				}
				return id, nil
			},
			base: func(id entities.PostID) (*http.Response, error) {
				return apiClient.Get(fmt.Sprintf("/posts/%v", id))
			},
			assertion: func(response *http.Response, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, response.StatusCode, http.StatusBadRequest)
					var message map[string]string
					if err := ParseJSONBody(response.Body, &message); assert.NoError(t, err) {
						if m, ok := message["message"]; assert.True(t, ok) {
							assert.Equal(t, m, entities.ErrPostNotFound.Error())
						}
					}
				}
			},
		},
	}

	for _, c := range cases {
		id, err := c.setup()
		assert.Nil(t, err)
		c.assertion(c.base(id))
	}
}
