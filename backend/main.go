package main

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize Elasticsearch client
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTICSEARCH_URL")},
	})
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
	}

	// Initialize repository and handler
	repo := NewSearchRepository(esClient)
	handler := NewSearchHandler(repo)

	// Initialize Fiber app
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return c.Next()
	})

	// Setup routes
	app.Get("/api/search", handler.Search)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
