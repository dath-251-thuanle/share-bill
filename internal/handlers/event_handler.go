package handlers

import (
	models "BACKEND/internal/dto"
	services "BACKEND/internal/services"
	utils "BACKEND/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(s *services.EventService) *EventHandler {
	return &EventHandler{service: s}
}


// CreateEvent POST /events
func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	var req models.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "INVALID_BODY", Message: "Invalid request body",
		})
	}

	resp, err := h.service.CreateEvent(c.Context(), userID, req)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true, Message: "Event created successfully", Data: resp,
	})
}

// GetEvent GET /events/:id
func (h *EventHandler) GetEvent(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	eventUUID := c.Params("id")

	resp, err := h.service.GetEvent(c.Context(), userID, eventUUID)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true, Data: resp,
	})
}

// UpdateEvent PUT /events/:id
func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	eventUUID := c.Params("id")
	var req models.UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "INVALID_BODY", Message: "Invalid request body",
		})
	}

	resp, err := h.service.UpdateEvent(c.Context(), userID, eventUUID, req)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true, Message: "Event updated successfully", Data: resp,
	})
}

// ListEvents GET /events
func (h *EventHandler) ListEvents(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	resp, err := h.service.ListEvents(c.Context(), userID)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true, Data: resp,
	})
}

// JoinEvent POST /events/:id/join
func (h *EventHandler) JoinEvent(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	eventUUID := c.Params("id")
	var req models.UpsertParticipantRequest
	_ = c.BodyParser(&req)

	err := h.service.JoinEvent(c.Context(), userID, eventUUID, req)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true, Message: "Joined event successfully",
	})
}

// LeaveEvent POST /events/:id/leave
func (h *EventHandler) LeaveEvent(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	eventUUID := c.Params("id")

	err := h.service.LeaveEvent(c.Context(), userID, eventUUID)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true, Message: "Left event successfully",
	})
}

// DeleteEvent DELETE /events/:id
func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	eventUUID := c.Params("id")

	err := h.service.DeleteEvent(c.Context(), userID, eventUUID)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true, Message: "Event deleted successfully",
	})
}