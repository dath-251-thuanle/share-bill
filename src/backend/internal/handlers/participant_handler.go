package handlers

import (
	"sharebill-backend/internal/models"
	"sharebill-backend/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ParticipantHandler struct {
	participantService *services.ParticipantService
}

func NewParticipantHandler(participantService *services.ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{
		participantService: participantService,
	}
}

func (h *ParticipantHandler) JoinEvent(c *fiber.Ctx) error {
	userIDstr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "INVALID USER ID"})
	}

	eventIDstr := strings.TrimSpace(c.Params("id"))
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "INVALID EVENT ID"})
	}
	err = h.participantService.JoinEvent(c.Context(), userID, eventID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "JOINED EVENT SUCCESSFULLY"})
}

func (h *ParticipantHandler) LeaveEvent(c *fiber.Ctx) error {
	userIDstr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "INVALID USER ID"})
	}

	eventIDstr := strings.TrimSpace(c.Params("id"))
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "INVALID EVENT ID"})
	}
	err = h.participantService.LeaveEvent(c.Context(), userID, eventID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "LEFT EVENT SUCCESSFULLY"})
}

func (h *ParticipantHandler) GetParticipants(c *fiber.Ctx) error {
	eventIDstr := c.Params("id")
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "INVALID EVENT ID"})
	}

	member, err := h.participantService.GetParticipants(c.Context(), eventID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "FETCH PARTICIPANTS SUCCESSFULLY",
		"participants": member,
	})
}

func (h *ParticipantHandler) GetMyEvents(c *fiber.Ctx) error {
	userIDstr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil{
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	events, err := h.participantService.GetMyParticipatedEvents(c.Context(), userID)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if events == nil{
		events = []*models.Event{}
	}

	return c.JSON(fiber.Map{
		"message": "Get list participated events successfully",
		"events": events,
		"total": len(events),
	})
}