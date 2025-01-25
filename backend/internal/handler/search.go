package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
	"github.com/oSoloTurk/multiple-kind-search/internal/logger"
)

type SearchHandler struct {
	searchService domain.SearchService
}

func NewSearchHandler(searchService domain.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

// Search godoc
// @Summary Search news with author boosting
// @Description Search news content with boosted results for specified author
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param username query string true "Author username to boost results for"
// @Success 200 {array} domain.SearchResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/search [get]
func (h *SearchHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		logger.Logger.Error().Msg("Search query is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	username := c.Query("username")
	if username == "" {
		logger.Logger.Error().Msg("Username is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}

	logger.Logger.Info().
		Str("query", query).
		Str("username", username).
		Msg("Processing search request")

	results, err := h.searchService.Search(c.UserContext(), domain.SearchFilter{
		Query:    query,
		Username: username,
	})
	if err != nil {
		logger.Logger.Error().
			Err(err).
			Str("query", query).
			Str("username", username).
			Msg("Failed to search news")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Logger.Info().
		Str("query", query).
		Str("username", username).
		Int("results", len(results)).
		Msg("Search completed successfully")

	return c.JSON(results)
}
