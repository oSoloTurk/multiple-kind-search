package domain

import "context"

type SearchResultType string

const (
	NewsResultType   SearchResultType = "news"
	AuthorResultType SearchResultType = "author"
)

type SearchResult struct {
	ID      string           `json:"id"`
	Title   string           `json:"title"`
	Content string           `json:"content"`
	Score   float64          `json:"score"`
	Type    SearchResultType `json:"type"`
}

type SearchFilter struct {
	Query    string
	Username string
}

type SearchService interface {
	Search(ctx context.Context, filter SearchFilter) ([]SearchResult, error)
}

type SearchRepository interface {
	Search(ctx context.Context, filter SearchFilter) ([]SearchResult, error)
}
