package service

import (
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
)

type newsService struct {
	repo domain.NewsRepository
}

func NewNewsService(repo domain.NewsRepository) domain.NewsService {
	return &newsService{repo: repo}
}

func (s *newsService) Create(news *domain.News) error {
	if err := news.Validate(); err != nil {
		return err
	}
	return s.repo.Create(news)
}

func (s *newsService) GetByID(id string) (*domain.News, error) {
	return s.repo.GetByID(id)
}

func (s *newsService) Update(news *domain.News) error {
	if err := news.Validate(); err != nil {
		return err
	}
	return s.repo.Update(news)
}

func (s *newsService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *newsService) List() ([]domain.News, error) {
	return s.repo.List()
}

func (s *newsService) Search(query string, username string) ([]domain.News, error) {
	return s.repo.Search(query, username)
}
