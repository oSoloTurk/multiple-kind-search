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
	log := logger.Logger.With().Str("query", query).Str("username", username).Logger()

	log.Info().Msg("Starting search operation")

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
			log = log.With().Str("authorID", authorID).Logger()
			log.Info().Msg("Found author for boosting")
		}
	}

	// Build the search query
	multiMatchQuery := elastic.NewMultiMatchQuery(query, "title", "content").
		Type("best_fields").
		TieBreaker(0.3)

	// For ES 7.10.2, we use bool query with should clauses instead of function_score
	boolQuery := elastic.NewBoolQuery().
		Must(multiMatchQuery)

	if authorID != "" {
		// Add author boost using should clause with boost parameter
		authorBoostQuery := elastic.NewTermQuery("authorID", authorID).Boost(2.0)
		boolQuery.Should(authorBoostQuery)
		log.Info().Msg("Applied author boost to search query")
	}

	querySource, err := boolQuery.Source()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get query source")
	} else {
		queryJSON, err := json.Marshal(querySource)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal query to JSON")
		} else {
			log.Info().RawJSON("query", queryJSON).Msg("Search query details")
		}
	}

	result, err := r.client.Search().
		Index(newsIndex).
		Query(boolQuery).
		Size(1000).
		Do(context.Background())

	if err != nil {
		log.Error().Err(err).Msg("Failed to execute search")
		return nil, err
	}

	newsList := make([]domain.News, 0)
	for _, hit := range result.Hits.Hits {
		var news domain.News
		if err := json.Unmarshal(hit.Source, &news); err != nil {
			log.Error().
				Err(err).
				Str("id", hit.Id).
				Msg("Failed to unmarshal news item")
			return nil, err
		}
		newsList = append(newsList, news)
	}

	log.Info().
		Int("resultCount", len(newsList)).
		Msg("Search completed successfully")

	return newsList, nil
}
