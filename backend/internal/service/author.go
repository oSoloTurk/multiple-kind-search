package service

import (
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
)

type authorService struct {
	repo domain.AuthorRepository
}

func NewAuthorService(repo domain.AuthorRepository) domain.AuthorService {
	return &authorService{repo: repo}
}

func (s *authorService) Create(author *domain.Author) error {
	if err := author.Validate(); err != nil {
		return err
	}
	return s.repo.Create(author)
}

func (s *authorService) GetByID(id string) (*domain.Author, error) {
	return s.repo.GetByID(id)
}

func (s *authorService) Update(author *domain.Author) error {
	if err := author.Validate(); err != nil {
		return err
	}
	return s.repo.Update(author)
}

func (s *authorService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *authorService) List() ([]domain.Author, error) {
	return s.repo.List()
}
