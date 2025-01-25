package elasticsearch

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
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

	_, err := r.client.Index().
		Index(newsIndex).
		Id(news.ID).
		BodyJson(news).
		Do(context.Background())

	return err
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
		}
	}

	// Build the search query
	multiMatchQuery := elastic.NewMultiMatchQuery(query, "title", "content").
		Type("best_fields").
		TieBreaker(0.3)

	// Create a function score query to boost author's content
	functionScoreQuery := elastic.NewFunctionScoreQuery().
		Query(multiMatchQuery)

	if authorID != "" {
		// Boost score by 2.0 if the author matches
		functionScoreQuery.Add(
			elastic.NewTermQuery("author_id", authorID),
			elastic.NewWeightFactorFunction(2.0),
		)
	}

	// Execute the search
	result, err := r.client.Search().
		Index(newsIndex).
		Query(functionScoreQuery).
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
