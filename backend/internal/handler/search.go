package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oSoloTurk/multiple-kind-search/internal/domain"
)

type SearchHandler struct {
	newsService domain.NewsService
}

func NewSearchHandler(newsService domain.NewsService) *SearchHandler {
	return &SearchHandler{newsService: newsService}
}

// Search godoc
// @Summary Search news with author boosting
// @Description Search news content with boosted results for specified author
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param username query string true "Author username to boost results for"
// @Success 200 {array} domain.News
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/search [get]
func (h *SearchHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	username := c.Query("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}

	news, err := h.newsService.Search(query, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(news)
}
