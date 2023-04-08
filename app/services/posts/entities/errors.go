package entities

import (
	"errors"
)

var (
	// ErrNilDB is thrown when an unexpected nil db connection is encountered.
	ErrNilDB = errors.New("db connection is nil")
	// ErrPostNotFound is thrown when a post is fetched for an invalid id.
	ErrPostNotFound = errors.New("no post found with provided id")
	// ErrDuplicatePost is thrown when a post with conflicting id is created.
	ErrDuplicatePost = errors.New("post id already exists")
	// ErrInvalidTitle is thrown when title is empty or longer than 255 characters.
	ErrInvalidTitle = errors.New("title has to be a non-empty string of 255 characters or less")
)
