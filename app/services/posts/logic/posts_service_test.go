package logic

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/database"
	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/entities"
	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/interfaces"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	postsServiceTestInstance interfaces.PostsService
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	teardown, err := setup(ctx)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to setup unit tests")
		os.Exit(1)
	}
	exitValue := m.Run()
	teardown(ctx)
	os.Exit(exitValue)
}

func TestNewPostsService(t *testing.T) {
	ctx := context.Background()
	reader, err := database.GetReader(ctx)
	if !assert.NoError(t, err) {
		return
	}
	writer, err := database.GetWriter(ctx)
	if !assert.NoError(t, err) {
		return
	}
	_, err = NewPostsService(ctx, reader, writer)
	assert.Nil(t, err)
	_, err = NewPostsService(ctx, nil, nil)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, entities.ErrNilDB)
	}
}

func TestPostsService_IndexPosts(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func() ([]entities.Post, error)
		assertion func(posts []entities.Post, err error)
		cleanup   func(entities.PostID) error
	}{
		// Non-empty table.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post.ID, err
			},
			base: func() ([]entities.Post, error) {
				return postsServiceTestInstance.IndexPosts(ctx)
			},
			assertion: func(posts []entities.Post, err error) {
				if assert.NoError(t, err) {
					assert.NotEmpty(t, posts)
				}
			},
			cleanup: func(id entities.PostID) error {
				return postsServiceTestInstance.DeletePost(ctx, id)
			},
		},
	}
	for _, c := range cases {
		postID, err := c.setup()
		assert.Nil(t, err)
		posts, err := c.base()
		c.assertion(posts, err)
		err = c.cleanup(postID)
		assert.Nil(t, err)
	}
}

func TestPostsService_FindPost(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(id entities.PostID) (entities.Post, error)
		assertion func(post entities.Post, err error)
		cleanup   func(entities.PostID) error
	}{
		// Existing record.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				return postsServiceTestInstance.FindPost(ctx, id)
			},
			assertion: func(post entities.Post, err error) {
				if assert.NoError(t, err) {
					assert.NotZero(t, post)
				}
			},
			cleanup: func(id entities.PostID) error {
				return postsServiceTestInstance.DeletePost(ctx, id)
			},
		},
		// Non-existing post.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				_ = postsServiceTestInstance.DeletePost(ctx, post.ID)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				return postsServiceTestInstance.FindPost(ctx, id)
			},
			assertion: func(post entities.Post, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, entities.ErrPostNotFound)
					assert.Zero(t, post)
				}
			},
			cleanup: func(id entities.PostID) error { return nil },
		},
	}
	for _, c := range cases {
		postID, err := c.setup()
		assert.Nil(t, err)
		post, err := c.base(postID)
		c.assertion(post, err)
		err = c.cleanup(postID)
		assert.Nil(t, err)
	}
}

func TestPostsService_UpdatePost(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(id entities.PostID) (entities.Post, error)
		assertion func(post entities.Post, err error)
		cleanup   func(entities.PostID) error
	}{
		// Existing record.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					ID:      id,
					Title:   "test-title-new",
					Content: "test-content-new",
				}
				err := postsServiceTestInstance.UpdatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.NoError(t, err) {
					assert.NotZero(t, post)
				}
			},
			cleanup: func(id entities.PostID) error { return nil },
		},
		// Non-existing post.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				_ = postsServiceTestInstance.DeletePost(ctx, post.ID)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					ID:      id,
					Title:   strings.Repeat("x", 255),
					Content: "test-content-new",
				}
				err := postsServiceTestInstance.UpdatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, entities.ErrPostNotFound)
				}
			},
			cleanup: func(id entities.PostID) error { return nil },
		},
	}
	for _, c := range cases {
		postID, err := c.setup()
		assert.Nil(t, err)
		c.assertion(c.base(postID))
		err = c.cleanup(postID)
		assert.Nil(t, err)
	}
}

func TestPostsService_CreatePost(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(id entities.PostID) (entities.Post, error)
		assertion func(post entities.Post, err error)
		cleanup   func(entities.PostID) error
	}{
		// Valid.
		{
			setup: func() (entities.PostID, error) { return 0, nil },
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.NoError(t, err) {
					assert.NotZero(t, post)
				}
			},
			cleanup: func(id entities.PostID) error {
				return postsServiceTestInstance.DeletePost(ctx, id)
			},
		},
		// Invalid title (empty).
		{
			setup: func() (entities.PostID, error) { return 0, nil },
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					Title:   "",
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, entities.ErrInvalidTitle)
				}
			},
			cleanup: func(id entities.PostID) error { return nil },
		},
		// Invalid title (overflow).
		{
			setup: func() (entities.PostID, error) { return 0, nil },
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 256),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, entities.ErrInvalidTitle)
				}
			},
			cleanup: func(id entities.PostID) error { return nil },
		},
		// Colliding id.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					ID:      id,
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, entities.ErrDuplicatePost)
				}
			},
			cleanup: func(id entities.PostID) error {
				return postsServiceTestInstance.DeletePost(ctx, id)
			},
		},
	}
	for _, c := range cases {
		postID, err := c.setup()
		assert.Nil(t, err)
		post, err := c.base(postID)
		c.assertion(post, err)
		if postID != 0 {
			err = c.cleanup(postID)
			assert.Nil(t, err)
			continue
		}
		if post.ID != 0 {
			err = c.cleanup(post.ID)
			assert.Nil(t, err)
			continue
		}
	}
}

func TestPostsService_DeletePost(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		setup     func() (entities.PostID, error)
		base      func(id entities.PostID) (entities.Post, error)
		assertion func(post entities.Post, err error)
		cleanup   func(entities.PostID) error
	}{
		// Valid.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					ID:      id,
					Title:   strings.Repeat("x", 255),
					Content: "test-content-updated",
				}
				err := postsServiceTestInstance.UpdatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.NoError(t, err) {
					assert.Equal(t, "test-content-updated", post.Content)
				}
			},
			cleanup: func(id entities.PostID) error {
				return postsServiceTestInstance.DeletePost(ctx, id)
			},
		},
		// Invalid title (empty).
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					ID:      id,
					Title:   "",
					Content: "test-content-updated",
				}
				err := postsServiceTestInstance.UpdatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, entities.ErrInvalidTitle)
				}
			},
			cleanup: func(id entities.PostID) error {
				return postsServiceTestInstance.DeletePost(ctx, id)
			},
		},
		// Invalid title (overflow).
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					ID:      id,
					Title:   strings.Repeat("x", 256),
					Content: "test-content-updated",
				}
				err := postsServiceTestInstance.UpdatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, entities.ErrInvalidTitle)
				}
			},
			cleanup: func(id entities.PostID) error {
				return postsServiceTestInstance.DeletePost(ctx, id)
			},
		},
		// Non-existing post.
		{
			setup: func() (entities.PostID, error) {
				post := entities.Post{
					Title:   strings.Repeat("x", 255),
					Content: "test-content",
				}
				err := postsServiceTestInstance.CreatePost(ctx, &post)
				_ = postsServiceTestInstance.DeletePost(ctx, post.ID)
				return post.ID, err
			},
			base: func(id entities.PostID) (entities.Post, error) {
				post := entities.Post{
					ID:      id,
					Title:   strings.Repeat("x", 255),
					Content: "test-content-updated",
				}
				err := postsServiceTestInstance.UpdatePost(ctx, &post)
				return post, err
			},
			assertion: func(post entities.Post, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, entities.ErrPostNotFound)
				}
			},
			cleanup: func(id entities.PostID) error { return nil },
		},
	}
	for _, c := range cases {
		postID, err := c.setup()
		assert.Nil(t, err)
		post, err := c.base(postID)
		c.assertion(post, err)
		err = c.cleanup(postID)
		assert.Nil(t, err)
	}
}

func setup(ctx context.Context) (func(context.Context), error) {
	reader, err := database.GetReader(ctx)
	if err != nil {
		return nil, err
	}
	writer, err := database.GetWriter(ctx)
	if err != nil {
		return nil, err
	}
	service, err := NewPostsService(ctx, reader, writer)
	if err != nil {
		return nil, err
	}
	postsServiceTestInstance = service
	return func(ctx context.Context) {}, nil
}
