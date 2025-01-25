package main

import (
	"os"

	_ "github.com/oSoloTurk/multiple-kind-search/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/olivere/elastic/v7"
)

func main() {
	// Initialize logger
	InitLogger()

	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(
		elastic.SetURL(os.Getenv("ELASTICSEARCH_URL")),
		elastic.SetSniff(false),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create Elasticsearch client")
	}
	logger.Info().Msg("Successfully connected to Elasticsearch")

	// Initialize repository and handlers
	repo := NewRepository(esClient)
	handler := NewHandler(repo)
	searchHandler := NewSearchHandler(repo.SearchRepository)

	// Initialize Fiber app
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	})

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:          "http://localhost:8080/swagger/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	// Setup routes
	api := app.Group("/api")
	api.Get("/search", searchHandler.Search)
	api.Get("/entries/:id", handler.GetEntry)
	api.Post("/entries", handler.CreateEntry)
	api.Put("/entries/:id", handler.UpdateEntry)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info().Str("port", port).Msg("Starting server")
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
