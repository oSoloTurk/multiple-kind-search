package config

import "os"

type Config struct {
	ElasticsearchURL string
	ServerPort       string
}

func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	return &Config{
		ElasticsearchURL: esURL,
		ServerPort:       port,
	}
}
