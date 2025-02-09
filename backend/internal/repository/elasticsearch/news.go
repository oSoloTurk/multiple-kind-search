package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
	"github.com/oSoloTurk/multiple-kind-search/internal/logger"
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

	body, err := json.Marshal(news)
	if err != nil {
		return err
	}

	_, err = r.client.Index(
		newsIndex,
		strings.NewReader(string(body)),
		r.client.Index.WithDocumentID(news.ID),
		r.client.Index.WithContext(context.Background()),
	)
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
	res, err := r.client.Get(
		newsIndex,
		id,
		r.client.Get.WithContext(context.Background()),
	)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to get news article: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	source, ok := result["_source"].(map[string]interface{})
	if !ok {
		return nil, nil // Document not found
	}

	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	var news domain.News
	if err := json.Unmarshal(sourceBytes, &news); err != nil {
		return nil, err
	}

	return &news, nil
}

func (r *newsRepository) Update(news *domain.News) error {
	news.UpdatedAt = time.Now()

	body, err := json.Marshal(map[string]interface{}{
		"doc": news,
	})
	if err != nil {
		return err
	}

	_, err = r.client.Update(
		newsIndex,
		news.ID,
		strings.NewReader(string(body)),
		r.client.Update.WithContext(context.Background()),
	)

	return err
}

func (r *newsRepository) Delete(id string) error {
	res, err := r.client.Delete(
		newsIndex,
		id,
		r.client.Delete.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (r *newsRepository) List() ([]domain.News, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": 1000,
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex(newsIndex),
		r.client.Search.WithBody(strings.NewReader(string(body))),
		r.client.Search.WithContext(context.Background()),
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
	newsList := make([]domain.News, 0)

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"]

		sourceBytes, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		var news domain.News
		if err := json.Unmarshal(sourceBytes, &news); err != nil {
			return nil, err
		}
		newsList = append(newsList, news)
	}

	return newsList, nil
}
