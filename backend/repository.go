package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

type SearchRepository struct {
	client *elastic.Client
}

type SearchResult struct {
	ID         string              `json:"id,omitempty"`
	Title      string              `json:"title,omitempty"`
	Content    string              `json:"content,omitempty"`
	Author     string              `json:"author,omitempty"`
	Type       string              `json:"type,omitempty"`
	Highlights map[string][]string `json:"highlights,omitempty"`
}

type Repository struct {
	client *elastic.Client
	*SearchRepository
}

type Entry struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func NewSearchRepository(client *elastic.Client) *SearchRepository {
	return &SearchRepository{client: client}
}

func NewRepository(client *elastic.Client) *Repository {
	return &Repository{
		client:           client,
		SearchRepository: NewSearchRepository(client),
	}
}

func (r *SearchRepository) Search(query string) ([]SearchResult, error) {
	logger.Info().Str("query", query).Msg("Starting search across all indices")
	var results []SearchResult

	authorResults, err := r.searchAuthors(query)
	if err != nil {
		logger.Error().Err(err).Str("query", query).Msg("Error searching authors")
		return nil, fmt.Errorf("error searching authors: %w", err)
	}
	logger.Debug().Int("count", len(authorResults)).Msg("Author search results")
	results = append(results, authorResults...)

	titleResults, err := r.searchTitles(query)
	if err != nil {
		logger.Error().Err(err).Str("query", query).Msg("Error searching titles")
		return nil, fmt.Errorf("error searching titles: %w", err)
	}
	logger.Debug().Int("count", len(titleResults)).Msg("Title search results")
	results = append(results, titleResults...)

	contentResults, err := r.searchContents(query)
	if err != nil {
		logger.Error().Err(err).Str("query", query).Msg("Error searching contents")
		return nil, fmt.Errorf("error searching contents: %w", err)
	}
	logger.Debug().Int("count", len(contentResults)).Msg("Content search results")
	results = append(results, contentResults...)

	logger.Info().
		Str("query", query).
		Int("total_results", len(results)).
		Msg("Search completed successfully")
	return results, nil
}

func (r *SearchRepository) searchAuthors(query string) ([]SearchResult, error) {
	var results []SearchResult

	searchQuery := elastic.NewBoolQuery().
		Should(
			elastic.NewMatchQuery("name", query).Boost(2.0),
			elastic.NewMatchQuery("bio", query),
		)

	highlighter := elastic.NewHighlight().
		Field("name").NumOfFragments(0).
		Field("bio").NumOfFragments(1).FragmentSize(150).
		PreTags("<em>").PostTags("</em>")

	searchResult, err := r.client.Search().
		Index("authors").
		Query(searchQuery).
		Highlight(highlighter).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var source map[string]interface{}
		err := json.Unmarshal(hit.Source, &source)
		if err != nil {
			continue
		}

		title := source["name"].(string)
		content := source["bio"].(string)

		if hit.Highlight != nil {
			if len(hit.Highlight["name"]) > 0 {
				title = hit.Highlight["name"][0]
			}
			if len(hit.Highlight["bio"]) > 0 {
				content = hit.Highlight["bio"][0]
			}
		}

		results = append(results, SearchResult{
			ID:      hit.Id,
			Title:   title,
			Content: content,
			Type:    "author",
		})
	}

	return results, nil
}

func (r *SearchRepository) searchTitles(query string) ([]SearchResult, error) {
	var results []SearchResult

	searchQuery := elastic.NewMatchQuery("title", query).
		Operator("OR").
		Fuzziness("AUTO")

	highlighter := elastic.NewHighlight().
		Field("title").NumOfFragments(0).
		PreTags("<em>").PostTags("</em>")

	searchResult, err := r.client.Search().
		Index("titles").
		Query(searchQuery).
		Highlight(highlighter).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var source map[string]interface{}
		err := json.Unmarshal(hit.Source, &source)
		if err != nil {
			continue
		}

		title := source["title"].(string)

		if hit.Highlight != nil && len(hit.Highlight["title"]) > 0 {
			title = hit.Highlight["title"][0]
		}

		results = append(results, SearchResult{
			ID:    hit.Id,
			Title: title,
			Type:  "news_title",
		})
	}

	return results, nil
}

func (r *SearchRepository) searchContents(query string) ([]SearchResult, error) {
	var results []SearchResult

	searchQuery := elastic.NewMatchQuery("content", query).
		Operator("OR").
		MinimumShouldMatch("75%")

	highlighter := elastic.NewHighlight().
		Field("content").NumOfFragments(1).FragmentSize(150).
		PreTags("<em>").PostTags("</em>")

	searchResult, err := r.client.Search().
		Index("contents").
		Query(searchQuery).
		Highlight(highlighter).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var source map[string]interface{}
		err := json.Unmarshal(hit.Source, &source)
		if err != nil {
			continue
		}

		content := source["content"].(string)

		if hit.Highlight != nil && len(hit.Highlight["content"]) > 0 {
			content = hit.Highlight["content"][0]
		}

		results = append(results, SearchResult{
			ID:      hit.Id,
			Content: content,
			Type:    "news_content",
		})
	}

	return results, nil
}

func (r *Repository) GetEntry(id string) (*Entry, error) {
	result, err := r.client.Get().
		Index("entries").
		Id(id).
		Do(context.Background())

	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, fmt.Errorf("entry not found")
		}
		return nil, err
	}

	var entry Entry
	err = json.Unmarshal(result.Source, &entry)
	if err != nil {
		return nil, err
	}
	entry.ID = result.Id

	return &entry, nil
}

func (r *Repository) CreateEntry(entry *Entry) error {
	if entry.ID == "" {
		entry.ID = generateID()
	}

	_, err := r.client.Index().
		Index("entries").
		Id(entry.ID).
		BodyJson(entry).
		Refresh("true").
		Do(context.Background())

	return err
}

func (r *Repository) UpdateEntry(entry *Entry) error {
	_, err := r.client.Index().
		Index("entries").
		Id(entry.ID).
		BodyJson(entry).
		Refresh("true").
		Do(context.Background())

	return err
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
