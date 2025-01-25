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
