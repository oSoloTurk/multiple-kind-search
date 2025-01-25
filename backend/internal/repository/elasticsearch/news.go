package elasticsearch

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
	"github.com/oSoloTurk/multiple-kind-search/internal/logger"
	"github.com/olivere/elastic/v7"
)

const newsIndex = "news"

type newsRepository struct {
	client *elastic.Client
}

func NewNewsRepository(client *elastic.Client) domain.NewsRepository {
	return &newsRepository{client: client}
}

func (r *newsRepository) Create(news *domain.News) error {
	if news.ID == "" {
		news.ID = uuid.New().String()
	}
	now := time.Now()
	news.CreatedAt = now
	news.UpdatedAt = now

	logger.Logger.Info().
		Str("id", news.ID).
		Str("title", news.Title).
		Msg("Creating new news article")

	_, err := r.client.Index().
		Index(newsIndex).
		Id(news.ID).
		BodyJson(news).
		Do(context.Background())

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Str("id", news.ID).
			Msg("Failed to create news article")
		return err
	}

	return nil
}

func (r *newsRepository) GetByID(id string) (*domain.News, error) {
	result, err := r.client.Get().
		Index(newsIndex).
		Id(id).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	var news domain.News
	if err := json.Unmarshal(result.Source, &news); err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *newsRepository) Update(news *domain.News) error {
	news.UpdatedAt = time.Now()

	_, err := r.client.Update().
		Index(newsIndex).
		Id(news.ID).
		Doc(news).
		Do(context.Background())

	return err
}

func (r *newsRepository) Delete(id string) error {
	_, err := r.client.Delete().
		Index(newsIndex).
		Id(id).
		Do(context.Background())

	return err
}

func (r *newsRepository) List() ([]domain.News, error) {
	result, err := r.client.Search().
		Index(newsIndex).
		Size(1000).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	newsList := make([]domain.News, 0)
	for _, hit := range result.Hits.Hits {
		var news domain.News
		if err := json.Unmarshal(hit.Source, &news); err != nil {
			return nil, err
		}
		newsList = append(newsList, news)
	}

	return newsList, nil
}

func (r *newsRepository) Search(query string, username string) ([]domain.News, error) {
	logger.Logger.Info().
		Str("query", query).
		Str("username", username).
		Msg("Starting search operation")

	// Get author by username first
	authorResult, err := r.client.Search().
		Index("authors").
		Query(elastic.NewMatchQuery("name", username)).
		Size(1).
		Do(context.Background())

	var authorID string
	if err == nil && len(authorResult.Hits.Hits) > 0 {
		var author struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(authorResult.Hits.Hits[0].Source, &author); err == nil {
			authorID = author.ID
			logger.Logger.Debug().
				Str("authorID", authorID).
				Str("username", username).
				Msg("Found author for boosting")
		}
	}

	// Build the search query
	multiMatchQuery := elastic.NewMultiMatchQuery(query, "title", "content").
		Type("best_fields").
		TieBreaker(0.3)

	functionScoreQuery := elastic.NewFunctionScoreQuery().
		Query(multiMatchQuery)

	if authorID != "" {
		functionScoreQuery.Add(
			elastic.NewTermQuery("author_id", authorID),
			elastic.NewWeightFactorFunction(2.0),
		)
		logger.Logger.Debug().
			Str("authorID", authorID).
			Msg("Applied author boost to search query")
	}

	result, err := r.client.Search().
		Index(newsIndex).
		Query(functionScoreQuery).
		Size(1000).
		Do(context.Background())

	if err != nil {
		logger.Logger.Error().
			Err(err).
			Str("query", query).
			Str("username", username).
			Msg("Failed to execute search")
		return nil, err
	}

	newsList := make([]domain.News, 0)
	for _, hit := range result.Hits.Hits {
		var news domain.News
		if err := json.Unmarshal(hit.Source, &news); err != nil {
			logger.Logger.Error().
				Err(err).
				Str("id", hit.Id).
				Msg("Failed to unmarshal news item")
			return nil, err
		}
		newsList = append(newsList, news)
	}

	logger.Logger.Info().
		Int("resultCount", len(newsList)).
		Str("query", query).
		Str("username", username).
		Msg("Search completed successfully")

	return newsList, nil
}
