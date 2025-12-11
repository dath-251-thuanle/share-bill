package services

import (
	"context"
	"errors"
	"sharebill-backend/internal/models"
	"sharebill-backend/internal/repositories"

	"github.com/google/uuid"
)

type EventService struct {
	eventRepo repositories.EventRepository
}

func NewEventService(eventRepo repositories.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

func (s *EventService) CreateEvent(ctx context.Context, userID uuid.UUID , name string, description *string) (*models.Event, error) {
	event := &models.Event{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		OwnerID:     userID,
	}
	err := s.eventRepo.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) GetMyEvents(ctx context.Context, userID uuid.UUID) ([]*models.Event, error) {
	return s.eventRepo.ListByUserID(ctx, userID)
}

func (s *EventService) GetEventByID(ctx context.Context, eventID uuid.UUID) (*models.EventDetailResponse, error) {
	return s.eventRepo.GetByID(ctx, eventID)
}

func (s *EventService) UpdateEvent(ctx context.Context, userID, eventID uuid.UUID, name string, description *string) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("EVENT NOT FOUND")
	}
	if event.OwnerID != userID {
		return errors.New("YOU ARE NOT AUTHORIZED TO UPDATE THIS EVENT")
	}

	return s.eventRepo.Update(ctx, eventID, name, description)
}

func (s *EventService) SettleEvent(ctx context.Context, userID, eventID uuid.UUID, isSettled bool) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("EVENT NOT FOUND")
	}
	if event.OwnerID != userID {
		return errors.New("YOU ARE NOT AUTHORIZED TO SETTLE THIS EVENT")
	}

	return s.eventRepo.SettleEvent(ctx, eventID, isSettled)
}

func (s *EventService) DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("EVENT NOT FOUND")
	}
	if event.OwnerID != userID {
		return errors.New("YOU ARE NOT AUTHORIZED TO DELETE THIS EVENT")
	}
	return s.eventRepo.Delete(ctx, eventID)
}