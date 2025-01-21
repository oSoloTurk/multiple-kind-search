package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

type SearchRepository struct {
	client *elasticsearch.Client
}

type SearchResult struct {
	Title      string              `json:"title,omitempty"`
	Content    string              `json:"content,omitempty"`
	Author     string              `json:"author,omitempty"`
	Type       string              `json:"type,omitempty"`
	Highlights map[string][]string `json:"highlights,omitempty"`
}

type Suggestion struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func NewSearchRepository(client *elasticsearch.Client) *SearchRepository {
	return &SearchRepository{client: client}
}

func (r *SearchRepository) Search(query string) ([]SearchResult, error) {
	var results []SearchResult

	// Search in all indices
	authorResults, err := r.searchAuthors(query)
	if err != nil {
		return nil, fmt.Errorf("error searching authors: %w", err)
	}
	results = append(results, authorResults...)

	titleResults, err := r.searchTitles(query)
	if err != nil {
		return nil, fmt.Errorf("error searching titles: %w", err)
	}
	results = append(results, titleResults...)

	contentResults, err := r.searchContents(query)
	if err != nil {
		return nil, fmt.Errorf("error searching contents: %w", err)
	}
	results = append(results, contentResults...)

	return results, nil
}

func (r *SearchRepository) searchAuthors(query string) ([]SearchResult, error) {
	var results []SearchResult

	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "bio"},
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

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("authors"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if hits, ok := response["hits"].(map[string]interface{}); ok {
		if hitList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitList {
				hitMap := hit.(map[string]interface{})
				source := hitMap["_source"].(map[string]interface{})

				title := source["name"].(string)
				content := source["bio"].(string)

				if highlightSection, ok := hitMap["highlight"].(map[string]interface{}); ok {
					if nameHighlights, ok := highlightSection["name"].([]interface{}); ok && len(nameHighlights) > 0 {
						title = nameHighlights[0].(string)
					}
					if bioHighlights, ok := highlightSection["bio"].([]interface{}); ok && len(bioHighlights) > 0 {
						content = bioHighlights[0].(string)
					}
				}

				results = append(results, SearchResult{
					Title:   title,
					Content: content,
					Type:    "author",
				})
			}
		}
	}

	return results, nil
}

func (r *SearchRepository) searchTitles(query string) ([]SearchResult, error) {
	var results []SearchResult

	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": query,
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
			},
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("titles"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if hits, ok := response["hits"].(map[string]interface{}); ok {
		if hitList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitList {
				hitMap := hit.(map[string]interface{})
				source := hitMap["_source"].(map[string]interface{})

				title := source["title"].(string)

				if highlightSection, ok := hitMap["highlight"].(map[string]interface{}); ok {
					if titleHighlights, ok := highlightSection["title"].([]interface{}); ok && len(titleHighlights) > 0 {
						title = titleHighlights[0].(string)
					}
				}

				results = append(results, SearchResult{
					Title: title,
					Type:  "news_title",
				})
			}
		}
	}

	return results, nil
}

func (r *SearchRepository) searchContents(query string) ([]SearchResult, error) {
	var results []SearchResult

	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": query,
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"content": map[string]interface{}{},
			},
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchBody); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("contents"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	if hits, ok := response["hits"].(map[string]interface{}); ok {
		if hitList, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitList {
				hitMap := hit.(map[string]interface{})
				source := hitMap["_source"].(map[string]interface{})

				content := source["content"].(string)

				if highlightSection, ok := hitMap["highlight"].(map[string]interface{}); ok {
					if contentHighlights, ok := highlightSection["content"].([]interface{}); ok && len(contentHighlights) > 0 {
						content = contentHighlights[0].(string)
					}
				}

				results = append(results, SearchResult{
					Content: content,
					Type:    "news_content",
				})
			}
		}
	}

	return results, nil
}
