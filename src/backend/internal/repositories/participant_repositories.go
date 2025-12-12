package repositories

import (
	"context"
	"errors"
	"sharebill-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ParticipantRepository interface {
	AddParticipant(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) error
	GetParticipantsByEventID(ctx context.Context, eventID uuid.UUID) ([]models.ParticipantResponse, error)
	RemoveParticipant(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Event, error)
}

type participantRepository struct {
	db *pgxpool.Pool
}

func NewParticipantRepository(db *pgxpool.Pool) ParticipantRepository {
	return &participantRepository{db: db}
}

func (r *participantRepository) AddParticipant(ctx context.Context, eventID, userID uuid.UUID) error {
	query := `INSERT INTO participants (event_id, user_id) VALUES ($1, $2) ON CONFLICT (event_id, user_id) DO NOTHING`
	_, err := r.db.Exec(ctx, query, eventID, userID)
	return err
}

func (r *participantRepository) GetParticipantsByEventID(ctx context.Context, eventID uuid.UUID) ([]models.ParticipantResponse, error) {
	query := `
		SELECT u.id, u.name, u.email, u.bank_name, u.account_number, u.account_name, p.joined_at
		FROM participants p
		JOIN users u ON p.user_id = u.id
		WHERE p.event_id = $1
		ORDER BY p.joined_at ASC
	`
	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.ParticipantResponse
	for rows.Next() {
		var p models.ParticipantResponse
		err := rows.Scan(&p.UserID, &p.Name, &p.Email, &p.BankName, &p.AccountNumber, &p.AccountName, &p.JoinedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *participantRepository) RemoveParticipant(ctx context.Context, eventID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM participants WHERE event_id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, eventID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("YOU HAVEN'T PARTICIPATED OR THIS EVENT IS NOT EXIST")
	}
	return nil
}

func (r *participantRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Event, error) {
	query := `SELECT e.id, e.name, e.description, e.owner_id, e.created_at, e.updated_at
			 FROM events e LEFT JOIN participants p ON e.id = p.event_id
			 WHERE p.user_id = $1`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next(){
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.OwnerID, &event.CreatedAt, &event.UpdatedAt)
		if err != nil{
			return  nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}