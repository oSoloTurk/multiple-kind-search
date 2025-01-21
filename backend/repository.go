package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type SearchRepository struct {
	client *elasticsearch.Client
}

type SearchResult struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Author  string `json:"author,omitempty"`
	Type    string `json:"type,omitempty"`
}

type Suggestion struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func NewSearchRepository(client *elasticsearch.Client) *SearchRepository {
	return &SearchRepository{client: client}
}

func (r *SearchRepository) Search(query string) ([]SearchResult, error) {
	// Build the search query
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title^2", "content", "author"},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	// Search across all indices
	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("titles", "content", "authors"),
		r.client.Search.WithBody(&buf),
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
	searchResults := make([]SearchResult, 0)

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		index := hit.(map[string]interface{})["_index"].(string)

		searchResult := SearchResult{}
		switch index {
		case "titles":
			searchResult.Title = source["title"].(string)
			searchResult.Type = "title"
		case "content":
			searchResult.Content = source["content"].(string)
			searchResult.Type = "content"
		case "authors":
			searchResult.Author = source["author"].(string)
			searchResult.Type = "author"
		}
		searchResults = append(searchResults, searchResult)
	}

	return searchResults, nil
}

func (r *SearchRepository) Suggest(query string) ([]Suggestion, error) {
	// Build the suggestion query
	suggestQuery := map[string]interface{}{
		"suggest": map[string]interface{}{
			"text": query,
			"titles": map[string]interface{}{
				"completion": map[string]interface{}{
					"field": "title_suggest",
					"size":  5,
				},
			},
			"content": map[string]interface{}{
				"completion": map[string]interface{}{
					"field": "content_suggest",
					"size":  5,
				},
			},
			"authors": map[string]interface{}{
				"completion": map[string]interface{}{
					"field": "author_suggest",
					"size":  5,
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(suggestQuery); err != nil {
		return nil, err
	}

	// Search across all indices
	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("titles", "content", "authors"),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	suggestions := make([]Suggestion, 0)
	suggest := result["suggest"].(map[string]interface{})

	// Process suggestions from each type
	for suggestType, suggestResults := range suggest {
		if suggestResults, ok := suggestResults.([]interface{}); ok {
			for _, option := range suggestResults[0].(map[string]interface{})["options"].([]interface{}) {
				suggestions = append(suggestions, Suggestion{
					Text: option.(map[string]interface{})["text"].(string),
					Type: suggestType,
				})
			}
		}
	}

	return suggestions, nil
}

func (r *SearchRepository) CreateIndices() error {
	indices := []string{"titles", "content", "authors"}
	mappings := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "text",
				},
				"title_suggest": map[string]interface{}{
					"type": "completion",
				},
				"content": map[string]interface{}{
					"type": "text",
				},
				"content_suggest": map[string]interface{}{
					"type": "completion",
				},
				"author": map[string]interface{}{
					"type": "text",
				},
				"author_suggest": map[string]interface{}{
					"type": "completion",
				},
			},
		},
	}

	for _, index := range indices {
		res, err := r.client.Indices.Create(
			index,
			r.client.Indices.Create.WithBody(strings.NewReader(mappings)),
		)
		if err != nil {
			return err
		}
		res.Body.Close()
	}

	return nil
} 