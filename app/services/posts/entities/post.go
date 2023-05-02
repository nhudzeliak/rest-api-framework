package entities

import (
	"time"

	"github.com/pkg/errors"
)

// Post represents a newsletter post.
type Post struct {
	ID        PostID    `json:"id" gorm:"column:id; primary_key:yes"`
	Title     string    `json:"title" gorm:"column:title"`
	Content   string    `json:"content" gorm:"column:content"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

// TableName ...
func (Post) TableName() string {
	return "posts"
}

// Validate checks whether a given Post object is valid.
func (p Post) Validate() error {
	errs := make([]error, 0)

	if len(p.Title) == 0 || len(p.Title) > 255 {
		errs = append(errs, ErrInvalidTitle)
	}

	if len(errs) == 0 {
		return nil
	}
	err := errs[0]
	for _, e := range errs[1:] {
		err = errors.Wrap(err, e.Error())
	}
	return err
}
