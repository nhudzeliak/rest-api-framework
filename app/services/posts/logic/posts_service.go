package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/entities"
	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/posts/interfaces"

	"gorm.io/gorm"
)

// Verify that PostsService satisfies the interfaces.PostsService interface.
// This should throw a compilation error otherwise.
var _ interfaces.PostsService = (*PostsService)(nil)

// PostsService implements interfaces.PostsService.
type PostsService struct {
	reader *gorm.DB
	writer *gorm.DB
}

// NewPostsService instantiates a new PostsService.
func NewPostsService(ctx context.Context, reader *gorm.DB, writer *gorm.DB) (*PostsService, error) {
	if reader == nil || writer == nil {
		return nil, entities.ErrNilDB
	}
	return &PostsService{
		reader: reader,
		writer: writer,
	}, nil
}

// IndexPosts returns an array of all existing posts, throws entities.ErrPostNotFound if table is empty.
func (s *PostsService) IndexPosts(ctx context.Context) ([]entities.Post, error) {
	var posts []entities.Post
	err := s.reader.Find(&posts).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrPostNotFound
		}
		return nil, err
	}
	return posts, nil
}

// FindPost fetches a post by provided id, throws entities.ErrPostNotFound if id is invalid.
func (s *PostsService) FindPost(ctx context.Context, id entities.PostID) (entities.Post, error) {
	var post entities.Post
	err := s.reader.First(&post, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Post{}, entities.ErrPostNotFound
		}
		return entities.Post{}, err
	}
	return post, nil
}

// UpdatePost updates a post in persistent repository, throws entities.ErrPostNotFound if id is invalid.
func (s *PostsService) UpdatePost(ctx context.Context, post *entities.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	_, err := s.FindPost(ctx, post.ID)
	if err != nil {
		return err
	}
	return s.writer.Save(&post).Error
}

// CreatePost creates a post in persistent repository, throws entities.ErrDuplicatePost if id is conflicting.
func (s *PostsService) CreatePost(ctx context.Context, post *entities.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	err := s.writer.Create(&post).Error
	// Check for duplicate key error, didn't find a check in gorm :(
	if err != nil && strings.Contains(err.Error(), "SQLSTATE 23505") {
		err = entities.ErrDuplicatePost
	}
	return err
}

// DeletePost deletes a post from persistent repository, throws entities.ErrPostNotFound if id is invalid.
func (s *PostsService) DeletePost(ctx context.Context, id entities.PostID) error {
	post, err := s.FindPost(ctx, id)
	if err != nil {
		return err
	}
	return s.writer.Delete(&post).Error
}
