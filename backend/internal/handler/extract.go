package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"snaptask/backend/internal/model"
	"snaptask/backend/internal/service"
)

type Handler struct {
	gemini *service.GeminiService
	google *service.GoogleService
}

func New(gemini *service.GeminiService, google *service.GoogleService) *Handler {
	return &Handler{gemini: gemini, google: google}
}

func (h *Handler) Register(app *fiber.App) {
	health := func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	}
	app.Get("/health", health)
	app.Get("/healthz", health)
	app.Post("/extract", h.Extract)
	app.Post("/push", h.Push)
}

func (h *Handler) Extract(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "missing image form field")
	}
	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "cannot open uploaded image")
	}
	defer file.Close()

	items, err := h.gemini.Extract(c.Context(), file, fileHeader)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	return c.JSON(model.ExtractResponse{Items: items})
}

func (h *Handler) Push(c *fiber.Ctx) error {
	var req model.PushRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid push request")
	}
	if strings.TrimSpace(req.AccessToken) == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing Google access token")
	}
	if len(req.Items) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no items to push")
	}
	results := h.google.Push(c.Context(), req.AccessToken, req.Items)
	return c.JSON(model.PushResponse{Results: results})
}
