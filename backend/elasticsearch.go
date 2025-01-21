package main

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticsearchClient(url string) (*elasticsearch.Client, error) {
	config := elasticsearch.Config{
		Addresses: []string{url},
	}
	
	return elasticsearch.NewClient(config)
} 