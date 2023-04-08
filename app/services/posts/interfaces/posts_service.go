package interfaces

import (
	"context"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/entities"
)

// PostsService is an interface that is used by outside packages to interact with posts.
type PostsService interface {
	IndexPosts(ctx context.Context) ([]entities.Post, error)
	FindPost(ctx context.Context, id entities.PostID) (entities.Post, error)
	UpdatePost(ctx context.Context, post *entities.Post) error
	CreatePost(ctx context.Context, post *entities.Post) error
	DeletePost(ctx context.Context, id entities.PostID) error
}
