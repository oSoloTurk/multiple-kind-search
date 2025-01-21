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

func (h *SearchHandler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	results, err := h.repo.Search(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}

func (h *SearchHandler) Suggest(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	suggestions, err := h.repo.Suggest(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(suggestions)
} 