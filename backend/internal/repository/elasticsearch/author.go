package elasticsearch

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
)

const authorIndex = "authors"

type authorRepository struct {
	client *es.Client
}

func NewAuthorRepository(client *es.Client) domain.AuthorRepository {
	return &authorRepository{client: client}
}

func (r *authorRepository) Create(author *domain.Author) error {
	if author.ID == "" {
		author.ID = uuid.New().String()
	}
	now := time.Now()
	author.CreatedAt = now
	author.UpdatedAt = now

	body, err := json.Marshal(author)
	if err != nil {
		return err
	}

	res, err := r.client.Index(
		authorIndex,
		strings.NewReader(string(body)),
		r.client.Index.WithDocumentID(author.ID),
		r.client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (r *authorRepository) GetByID(id string) (*domain.Author, error) {
	res, err := r.client.Get(
		authorIndex,
		id,
		r.client.Get.WithContext(context.Background()),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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

	var author domain.Author
	if err := json.Unmarshal(sourceBytes, &author); err != nil {
		return nil, err
	}

	return &author, nil
}

func (r *authorRepository) Update(author *domain.Author) error {
	author.UpdatedAt = time.Now()

	body, err := json.Marshal(map[string]interface{}{
		"doc": author,
	})
	if err != nil {
		return err
	}

	res, err := r.client.Update(
		authorIndex,
		author.ID,
		strings.NewReader(string(body)),
		r.client.Update.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (r *authorRepository) Delete(id string) error {
	res, err := r.client.Delete(
		authorIndex,
		id,
		r.client.Delete.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (r *authorRepository) List() ([]domain.Author, error) {
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
		r.client.Search.WithIndex(authorIndex),
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
	authors := make([]domain.Author, 0)

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"]

		sourceBytes, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		var author domain.Author
		if err := json.Unmarshal(sourceBytes, &author); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}
