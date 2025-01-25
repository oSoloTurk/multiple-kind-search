package main

import (
	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	repo *SearchRepository
}

func NewSearchHandler(repo *SearchRepository) *SearchHandler {
	return &SearchHandler{repo: repo}
}

// Search godoc
// @Summary Search entries
// @Description Search entries by query string
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} Entry
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/search [get]
func (h *SearchHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		logger.Warn().Msg("Empty search query received")
		return c.Status(400).JSON(map[string]string{
			"error": "Query parameter 'q' is required",
		})
	}

	logger.Info().Str("query", query).Msg("Processing search request")
	results, err := h.repo.Search(query)
	if err != nil {
		logger.Error().Err(err).Str("query", query).Msg("Search failed")
		return c.Status(500).JSON(map[string]string{
			"error": err.Error(),
		})
	}

	logger.Info().Str("query", query).Int("results", len(results)).Msg("Search completed")
	return c.JSON(results)
}

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	if repo == nil {
		logger.Fatal().Msg("Repository cannot be nil")
	}
	return &Handler{repo: repo}
}

// GetEntry godoc
// @Summary Get an entry by ID
// @Description Get a single entry by its ID
// @Tags entries
// @Accept json
// @Produce json
// @Param id path string true "Entry ID"
// @Success 200 {object} Entry
// @Failure 404 {object} map[string]string
// @Router /api/entries/{id} [get]
func (h *Handler) GetEntry(c *fiber.Ctx) error {
	id := c.Params("id")
	logger.Info().Str("id", id).Msg("Getting entry")

	entry, err := h.repo.GetEntry(id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Entry not found")
		return c.Status(404).JSON(map[string]string{
			"error": "Entry not found",
		})
	}

	logger.Info().Str("id", id).Msg("Entry retrieved successfully")
	return c.JSON(entry)
}

// CreateEntry godoc
// @Summary Create a new entry
// @Description Create a new entry with the provided data
// @Tags entries
// @Accept json
// @Produce json
// @Param entry body Entry true "Entry object"
// @Success 201 {object} Entry
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/entries [post]
func (h *Handler) CreateEntry(c *fiber.Ctx) error {
	var entry Entry
	if err := c.BodyParser(&entry); err != nil {
		logger.Error().Err(err).Msg("Invalid request body for create entry")
		return c.Status(400).JSON(map[string]string{
			"error": "Invalid request body",
		})
	}

	if entry.Title == "" || entry.Content == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "Title and content are required",
		})
	}

	logger.Info().Str("title", entry.Title).Msg("Creating new entry")
	err := h.repo.CreateEntry(&entry)
	if err != nil {
		logger.Error().Err(err).Interface("entry", entry).Msg("Failed to create entry")
		return c.Status(500).JSON(map[string]string{
			"error": "Database connection error: " + err.Error(),
		})
	}

	logger.Info().Str("id", entry.ID).Msg("Entry created successfully")
	return c.Status(201).JSON(entry)
}

// UpdateEntry godoc
// @Summary Update an existing entry
// @Description Update an entry by its ID with the provided data
// @Tags entries
// @Accept json
// @Produce json
// @Param id path string true "Entry ID"
// @Param entry body Entry true "Entry object"
// @Success 200 {object} Entry
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/entries/{id} [put]
func (h *Handler) UpdateEntry(c *fiber.Ctx) error {
	id := c.Params("id")
	var entry Entry
	if err := c.BodyParser(&entry); err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Invalid request body for update entry")
		return c.Status(400).JSON(map[string]string{
			"error": "Invalid request body",
		})
	}

	entry.ID = id
	logger.Info().Str("id", id).Msg("Updating entry")
	err := h.repo.UpdateEntry(&entry)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Interface("entry", entry).Msg("Failed to update entry")
		return c.Status(500).JSON(map[string]string{
			"error": "Failed to update entry",
		})
	}

	logger.Info().Str("id", id).Msg("Entry updated successfully")
	return c.JSON(entry)
}
