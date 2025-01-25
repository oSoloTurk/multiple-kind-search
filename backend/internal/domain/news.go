package domain

import (
	"errors"
	"time"
)

var (
	ErrNewsTitleRequired   = errors.New("news title is required")
	ErrNewsContentRequired = errors.New("news content is required")
	ErrNewsAuthorRequired  = errors.New("news author is required")
)

type News struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"authorID"`
	Tags      []string  `json:"tags,omitempty"`
	ImageURL  string    `json:"imageUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (n *News) Validate() error {
	if n.Title == "" {
		return ErrNewsTitleRequired
	}
	if n.Content == "" {
		return ErrNewsContentRequired
	}
	if n.AuthorID == "" {
		return ErrNewsAuthorRequired
	}
	return nil
}

type NewsRepository interface {
	Create(news *News) error
	GetByID(id string) (*News, error)
	Update(news *News) error
	Delete(id string) error
	List() ([]News, error)
}

type NewsService interface {
	Create(news *News) error
	GetByID(id string) (*News, error)
	Update(news *News) error
	Delete(id string) error
	List() ([]News, error)
}
