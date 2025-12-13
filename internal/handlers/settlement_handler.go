package handlers

import (
	models "BACKEND/internal/dto"
	services "BACKEND/internal/services"
	utils "BACKEND/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type SettlementHandler struct {
	service *services.SettlementService
}

func NewSettlementHandler(s *services.SettlementService) *SettlementHandler {
	return &SettlementHandler{service: s}
}

// GET /api/v1/events/:eventId/summary
func (h *SettlementHandler) GetEventSummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	eventUUID := c.Params("eventId")

	resp, err := h.service.GetEventSummary(c.Context(), userID, eventUUID)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// POST /api/v1/events/:eventId/settlements
func (h *SettlementHandler) CreateSettlement(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	eventUUID := c.Params("eventId")

	var req models.CreateSettlementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "INVALID_BODY", Message: "Invalid JSON format",
		})
	}

	err := h.service.CreateSettlement(c.Context(), userID, eventUUID, req)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true,
		Message: "Settlement recorded successfully",
	})
}