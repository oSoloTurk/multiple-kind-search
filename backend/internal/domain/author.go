package domain

import (
	"errors"
	"time"
)

var (
	ErrAuthorNameRequired = errors.New("author name is required")
)

type Author struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio,omitempty"`
	ImageURL  string    `json:"imageUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Author) Validate() error {
	if a.Name == "" {
		return ErrAuthorNameRequired
	}
	return nil
}

type AuthorRepository interface {
	Create(author *Author) error
	GetByID(id string) (*Author, error)
	Update(author *Author) error
	Delete(id string) error
	List() ([]Author, error)
}

type AuthorService interface {
	Create(author *Author) error
	GetByID(id string) (*Author, error)
	Update(author *Author) error
	Delete(id string) error
	List() ([]Author, error)
}
