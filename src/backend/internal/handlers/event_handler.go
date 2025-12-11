package handlers

import (
	"sharebill-backend/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	UserIDstr := c.Locals("user_id").(string)
	UserID, err := uuid.Parse(UserIDstr)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid token"})
	}

	var req struct {
		Name        string  `json:"name" validate:"required,min=3,max=100"`
		Description *string `json:"description,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	event, err := h.eventService.CreateEvent(c.Context(), UserID, req.Name, req.Description)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Event created successfully",
		"event":   event,
	})
}

func (h *EventHandler) GetMyEvents(c *fiber.Ctx) error {
	UserIDstr := c.Locals("user_id").(string)
	uid, _ := uuid.Parse(UserIDstr)

	events, err := h.eventService.GetMyEvents(c.Context(), uid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Events retrieved successfully",
		"events":  events,
	})
}

func (h* EventHandler) GetEventDetail(c *fiber.Ctx) error {
	eventIDstr := strings.TrimSpace(c.Params("id"))
	if eventIDstr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Event ID is required"})
	}
	
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		cleanID := strings.ReplaceAll(eventIDstr, " ", "")
		cleanID = strings.ToLower(cleanID)
		eventID, err = uuid.Parse(cleanID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid Event ID"})
		}
	}

	event, err := h.eventService.GetEventByID(c.Context(), eventID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if event == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Event not found"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Event retrieved successfully",
		"event":   event,
	})
}

func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	UserIDstr := c.Locals("user_id").(string)
	UserID, err := uuid.Parse(UserIDstr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	eventIDstr := c.Params("id")
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Event ID"})
	}

	var req struct {
		Name        string  `json:"name" validate:"required,min=3,max=100"`
		Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	err = h.eventService.UpdateEvent(c.Context(), UserID, eventID, req.Name, req.Description)
	if err != nil {
		if err.Error() == "You are not authorized to update this event" {
			return c.Status(403).JSON(fiber.Map{"error": "You are not authorized to update this event"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Event updated successfully",
	})
}

func (h *EventHandler) SettleEvent(c *fiber.Ctx) error {
	userIDstr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	eventIDstr := c.Params("id")
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Event ID"})
	}

	var req struct {
		IsSettled bool `json:"is_settled" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	err = h.eventService.SettleEvent(c.Context(), userID, eventID, req.IsSettled)
	if err != nil {
		if err.Error() == "You are not authorized to settle this event" {
			return c.Status(403).JSON(fiber.Map{"error": "You are not authorized to settle this event"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	action := "settled"
	if !req.IsSettled {
		action = "unsettled"
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Event " + action + " successfully",
	})
}

func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	userIDstr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	eventIDstr := c.Params("id")
	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Event ID"})
	}

	err = h.eventService.DeleteEvent(c.Context(), userID, eventID)
	if err != nil {
		if err.Error() == "You are not authorized to delete this event" {
			return c.Status(403).JSON(fiber.Map{"error": "You are not authorized to delete this event"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Event deleted successfully",
	})
}
