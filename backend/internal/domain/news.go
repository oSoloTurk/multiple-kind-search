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
	AuthorID  string    `json:"author_id"`
	Tags      []string  `json:"tags,omitempty"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Search(query string, username string) ([]News, error)
}

type NewsService interface {
	Create(news *News) error
	GetByID(id string) (*News, error)
	Update(news *News) error
	Delete(id string) error
	List() ([]News, error)
	Search(query string, username string) ([]News, error)
}
