package elasticsearch

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"sync"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
	"github.com/oSoloTurk/multiple-kind-search/internal/logger"
)

type SearchRepository struct {
	client *es.Client
}

func NewSearchRepository(client *es.Client) domain.SearchRepository {
	return &SearchRepository{client: client}
}

func (r *SearchRepository) Search(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, error) {
	results := make([]domain.SearchResult, 0)
	log := logger.Logger.With().Str("query", filter.Query).Str("username", filter.Username).Logger()
	log.Info().Msg("Starting combined search operation")

	// Create a wait group to wait for both goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Use slices to collect results
	var authorResults []domain.SearchResult
	var newsResults []domain.SearchResult

	// Search for authors asynchronously
	go func() {
		defer wg.Done()
		log.Info().Msg("Searching for authors")
		authors, err := r.SearchAuthor(ctx, filter)
		if err != nil {
			log.Error().Err(err).Msg("Error searching for authors")
			return
		}
		authorResults = authors
		log.Info().Int("authorCount", len(authors)).Msg("Authors search completed")
	}()

	// Search for news asynchronously
	go func() {
		defer wg.Done()
		log.Info().Msg("Searching for news")
		news, err := r.SearchNews(ctx, filter)
		if err != nil {
			log.Error().Err(err).Msg("Error searching for news")
			return
		}
		newsResults = news
		log.Info().Int("newsCount", len(news)).Msg("News search completed")
	}()

	wg.Wait()

	// Combine results
	results = append(results, authorResults...)
	results = append(results, newsResults...)

	// Sort results by score in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	log.Info().Int("totalResults", len(results)).Msg("Combined search operation completed successfully")
	return results, nil
}

func (r *SearchRepository) SearchAuthor(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, error) {
	// Build the search query for authors
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":       filter.Query,
				"fields":      []string{"name", "bio"},
				"type":        "best_fields",
				"tie_breaker": 0.3,
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"name": map[string]interface{}{},
				"bio":  map[string]interface{}{},
			},
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex("authors"),
		r.client.Search.WithBody(strings.NewReader(string(body))),
		r.client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	authors := make([]domain.SearchResult, 0)

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		highlights := hitMap["highlight"].(map[string]interface{})
		score := hitMap["_score"].(float64)

		var author domain.Author
		sourceBytes, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(sourceBytes, &author); err != nil {
			return nil, err
		}

		authors = append(authors, domain.SearchResult{
			ID:      author.ID,
			Title:   GetValueWithHighlight(highlights, "name", author.Name),
			Content: GetValueWithHighlight(highlights, "bio", author.Bio),
			Score:   score,
			Type:    domain.AuthorResultType,
		})
	}

	return authors, nil
}

func (r *SearchRepository) SearchNews(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, error) {
	// First find author ID if username is provided
	var authorID string
	if filter.Username != "" {
		authorQuery := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"name": filter.Username,
				},
			},
			"size": 1,
		}

		authorBody, err := json.Marshal(authorQuery)
		if err != nil {
			return nil, err
		}

		authorRes, err := r.client.Search(
			r.client.Search.WithIndex("authors"),
			r.client.Search.WithBody(strings.NewReader(string(authorBody))),
			r.client.Search.WithContext(ctx),
		)
		if err != nil {
			return nil, err
		}
		defer authorRes.Body.Close()

		var authorResult map[string]interface{}
		if err := json.NewDecoder(authorRes.Body).Decode(&authorResult); err != nil {
			return nil, err
		}

		hits := authorResult["hits"].(map[string]interface{})["hits"].([]interface{})
		if len(hits) > 0 {
			hitMap := hits[0].(map[string]interface{})
			source := hitMap["_source"].(map[string]interface{})
			authorID = source["id"].(string)
		}
	}

	// Build the search query for news
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":       filter.Query,
						"fields":      []string{"title", "content"},
						"type":        "best_fields",
						"tie_breaker": 0.3,
					},
				},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title":   map[string]interface{}{},
				"content": map[string]interface{}{},
			},
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
		},
	}

	// Add author boost if we have an author ID
	if authorID != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = map[string]interface{}{
			"term": map[string]interface{}{
				"authorID": map[string]interface{}{
					"value": authorID,
					"boost": 2.0,
				},
			},
		}
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex("news"),
		r.client.Search.WithBody(strings.NewReader(string(body))),
		r.client.Search.WithContext(ctx),
		r.client.Search.WithSize(1000),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	newsResults := make([]domain.SearchResult, 0)

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		highlights := hitMap["highlight"].(map[string]interface{})
		score := hitMap["_score"].(float64)

		var news domain.News
		sourceBytes, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(sourceBytes, &news); err != nil {
			return nil, err
		}

		newsResults = append(newsResults, domain.SearchResult{
			ID:      news.ID,
			Title:   GetValueWithHighlight(highlights, "title", news.Title),
			Content: GetValueWithHighlight(highlights, "content", news.Content),
			Score:   score,
			Type:    domain.NewsResultType,
		})
	}

	return newsResults, nil
}
