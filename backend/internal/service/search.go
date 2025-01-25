package service

import (
	"context"

	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
)

type SearchService struct {
	repo domain.SearchRepository
}

func NewSearchService(repo domain.SearchRepository) domain.SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) Search(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, error) {
	return s.repo.Search(ctx, filter)
}
