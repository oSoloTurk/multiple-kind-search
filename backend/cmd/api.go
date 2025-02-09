package cmd

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	elastic "github.com/elastic/go-elasticsearch/v8"

	"github.com/oSoloTurk/multiple-kind-search/internal/config"
	"github.com/oSoloTurk/multiple-kind-search/internal/handler"
	"github.com/oSoloTurk/multiple-kind-search/internal/logger"
	"github.com/oSoloTurk/multiple-kind-search/internal/repository/elasticsearch"
	"github.com/oSoloTurk/multiple-kind-search/internal/service"
	"github.com/spf13/cobra"
)

var (
	port string
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Long:  `Start the API server that handles all HTTP requests.`,
	Run:   runAPI,
}

func init() {
	apiCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server on")
	rootCmd.AddCommand(apiCmd)
}

func runAPI(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg := config.New()
	cfg.ServerPort = port // Override with flag value

	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(elastic.Config{
		Addresses: []string{cfg.ElasticsearchURL},
	})
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Initialize repositories
	authorRepo := elasticsearch.NewAuthorRepository(esClient)
	newsRepo := elasticsearch.NewNewsRepository(esClient)
	searchRepo := elasticsearch.NewSearchRepository(esClient)

	// Initialize services
	authorService := service.NewAuthorService(authorRepo)
	newsService := service.NewNewsService(newsRepo)
	searchService := service.NewSearchService(searchRepo)

	// Initialize handlers
	authorHandler := handler.NewAuthorHandler(authorService)
	newsHandler := handler.NewNewsHandler(newsService)
	searchHandler := handler.NewSearchHandler(searchService)

	// Initialize Fiber app
	app := fiber.New()

	// CORS middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	})

	// Swagger setup
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "http://localhost:" + cfg.ServerPort + "/swagger/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger")
	})

	// Routes
	api := app.Group("/api")

	// Search routes
	api.Get("/search", searchHandler.Search)

	// Domain routes
	authors := api.Group("/authors")
	authors.Post("/", authorHandler.Create)
	authors.Get("/", authorHandler.List)
	authors.Get("/:id", authorHandler.GetByID)
	authors.Put("/:id", authorHandler.Update)
	authors.Delete("/:id", authorHandler.Delete)

	news := api.Group("/news")
	news.Post("/", newsHandler.Create)
	news.Get("/", newsHandler.List)
	news.Get("/:id", newsHandler.GetByID)
	news.Put("/:id", newsHandler.Update)
	news.Delete("/:id", newsHandler.Delete)

	logger.Logger.Info().Msgf("Starting API on port %s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
