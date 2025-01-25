package elasticsearch

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
	"github.com/olivere/elastic/v7"
)

const authorIndex = "authors"

type authorRepository struct {
	client *elastic.Client
}

func NewAuthorRepository(client *elastic.Client) domain.AuthorRepository {
	return &authorRepository{client: client}
}

func (r *authorRepository) Create(author *domain.Author) error {
	if author.ID == "" {
		author.ID = uuid.New().String()
	}
	now := time.Now()
	author.CreatedAt = now
	author.UpdatedAt = now

	_, err := r.client.Index().
		Index(authorIndex).
		Id(author.ID).
		BodyJson(author).
		Do(context.Background())

	return err
}

func (r *authorRepository) GetByID(id string) (*domain.Author, error) {
	result, err := r.client.Get().
		Index(authorIndex).
		Id(id).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	var author domain.Author
	if err := json.Unmarshal(result.Source, &author); err != nil {
		return nil, err
	}

	return &author, nil
}

func (r *authorRepository) Update(author *domain.Author) error {
	author.UpdatedAt = time.Now()

	_, err := r.client.Update().
		Index(authorIndex).
		Id(author.ID).
		Doc(author).
		Do(context.Background())

	return err
}

func (r *authorRepository) Delete(id string) error {
	_, err := r.client.Delete().
		Index(authorIndex).
		Id(id).
		Do(context.Background())

	return err
}

func (r *authorRepository) List() ([]domain.Author, error) {
	result, err := r.client.Search().
		Index(authorIndex).
		Size(1000).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	authors := make([]domain.Author, 0)
	for _, hit := range result.Hits.Hits {
		var author domain.Author
		if err := json.Unmarshal(hit.Source, &author); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}
