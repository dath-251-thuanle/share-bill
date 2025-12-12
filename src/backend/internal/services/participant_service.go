package services

import (
	"context"
	"errors"
	"sharebill-backend/internal/models"
	"sharebill-backend/internal/repositories"

	"github.com/google/uuid"
)

type ParticipantService struct {
	participantRepo repositories.ParticipantRepository
	eventRepo repositories.EventRepository
}

func NewParticipantService(participantRepo repositories.ParticipantRepository, eventRepo repositories.EventRepository) *ParticipantService {
	return &ParticipantService{
		participantRepo: participantRepo,
		eventRepo: eventRepo,
	}
}

func (s *ParticipantService) JoinEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return errors.New("EVENT DOESN'T EXIST")
	}
	return s.participantRepo.AddParticipant(ctx, eventID, userID)
}

func (s *ParticipantService) LeaveEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil || event == nil {
		return errors.New("EVENT DOESN'T EXIST")
	}
	if event.OwnerID == userID {
		return errors.New("OWNER CANNOT LEAVE THE EVENT")
	}
	return s.participantRepo.RemoveParticipant(ctx, eventID, userID)
}

func (s *ParticipantService) GetParticipants(ctx context.Context, eventID uuid.UUID) ([]models.ParticipantResponse, error) {
	return s.participantRepo.GetParticipantsByEventID(ctx, eventID)
}

func (s *ParticipantService) GetMyParticipatedEvents(ctx context.Context, userID uuid.UUID) ([]*models.Event, error) {
	return s.participantRepo.ListByUserID(ctx, userID)
}