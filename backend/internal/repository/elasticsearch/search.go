package elasticsearch

import (
	"context"
	"encoding/json"
	"sort"
	"sync"

	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
	"github.com/oSoloTurk/multiple-kind-search/internal/logger"
	"github.com/olivere/elastic/v7"
)

type SearchRepository struct {
	client *elastic.Client
}

func NewSearchRepository(client *elastic.Client) domain.SearchRepository {
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
	multiMatchQuery := elastic.NewMultiMatchQuery(filter.Query, "name", "bio").
		Type("best_fields").
		TieBreaker(0.3)

	highlight := elastic.NewHighlight().
		Field("name").
		Field("bio").
		PreTags("<em>").
		PostTags("</em>")

	// Execute the search query
	searchResult, err := r.client.Search().
		Index("authors").
		Query(multiMatchQuery).
		Highlight(highlight).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	// Parse the search results
	authors := make([]domain.SearchResult, 0)
	for _, hit := range searchResult.Hits.Hits {
		var author domain.Author
		if err := json.Unmarshal(hit.Source, &author); err != nil {
			return nil, err
		}
		authors = append(authors, domain.SearchResult{
			ID:      author.ID,
			Title:   GetValueWithHighlight(hit.Highlight, "name", author.Name),
			Content: GetValueWithHighlight(hit.Highlight, "bio", author.Bio),
			Score:   *hit.Score,
			Type:    domain.AuthorResultType,
		})
	}

	return authors, nil
}

func (r *SearchRepository) SearchNews(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, error) {
	authorResult, err := r.client.Search().
		Index("authors").
		Query(elastic.NewMatchQuery("name", filter.Username)).
		Size(1).
		Do(ctx)

	var authorID string
	if err == nil && len(authorResult.Hits.Hits) > 0 {
		var author struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(authorResult.Hits.Hits[0].Source, &author); err == nil {
			authorID = author.ID
		}
	}

	// Build the search query
	multiMatchQuery := elastic.NewMultiMatchQuery(filter.Query, "title", "content").
		Type("best_fields").
		TieBreaker(0.3)

	// For ES 7.10.2, we use bool query with should clauses instead of function_score
	boolQuery := elastic.NewBoolQuery().
		Must(multiMatchQuery)

	if authorID != "" {
		// Add author boost using should clause with boost parameter
		authorBoostQuery := elastic.NewTermQuery("authorID", authorID).Boost(2.0)
		boolQuery.Should(authorBoostQuery)
	}

	highlight := elastic.NewHighlight().
		Field("title").
		Field("content").
		PreTags("<em>").
		PostTags("</em>")

	result, err := r.client.Search().
		Index("news").
		Query(boolQuery).
		Highlight(highlight).
		Size(1000).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	newsResults := make([]domain.SearchResult, 0)
	for _, hit := range result.Hits.Hits {
		var news domain.News
		if err := json.Unmarshal(hit.Source, &news); err != nil {
			return nil, err
		}
		newsResults = append(newsResults, domain.SearchResult{
			ID:      news.ID,
			Title:   GetValueWithHighlight(hit.Highlight, "title", news.Title),
			Content: GetValueWithHighlight(hit.Highlight, "content", news.Content),
			Score:   *hit.Score,
			Type:    domain.NewsResultType,
		})
	}

	return newsResults, nil
}
