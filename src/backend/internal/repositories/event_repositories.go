package repositories

import (
	"errors"
	"context"
	"sharebill-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Event, error)
	GetByID(ctx context.Context, eventID uuid.UUID) (*models.EventDetailResponse, error)
	Update(ctx context.Context, eventID uuid.UUID, name string, description *string) error
	SettleEvent(ctx context.Context, eventID uuid.UUID, isSettled bool) error
	Delete(ctx context.Context, eventID uuid.UUID) error
}

type eventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *models.Event) error {
	query := `INSERT INTO events (id, name, description, owner_id, is_settled)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		event.ID,
		event.Name,
		event.Description,
		event.OwnerID,
		false,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *eventRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Event, error) {
	query := `SELECT id ,name, description, owner_id, created_at, updated_at FROM events WHERE owner_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.OwnerID, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (r *eventRepository) GetByID(ctx context.Context, eventID uuid.UUID) (*models.EventDetailResponse, error) {
	query := `SELECT e.id, e.name, e.description, e.owner_id, e.created_at, e.updated_at, u.name, u.email FROM events e LEFT JOIN users u ON e.owner_id = u.id WHERE e.id = $1`
	var resp models.EventDetailResponse
	err := r.db.QueryRow(ctx, query, eventID).Scan(
		&resp.ID,
		&resp.Name,
		&resp.Description,
		&resp.OwnerID,
		&resp.CreatedAt,
		&resp.UpdatedAt,
		&resp.OwnerName,
		&resp.OwnerEmail,
	)
	if err == pgx.ErrNoRows{
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *eventRepository) Update(ctx context.Context, eventID uuid.UUID, name string, description *string) error {
	query := `UPDATE events SET name = $1, description = $2, updated_at = NOW() WHERE id = $3`
	
	result, err := r.db.Exec(ctx, query, name, description, eventID)
	if err != nil {
		return err
	}
	
	if result.RowsAffected() == 0 {
		return errors.New("YOU ARE NOT ALLOWED TO UPDATE THIS EVENT OR EVENT NOT FOUND")
	}
	return nil
}

func (r *eventRepository) SettleEvent(ctx context.Context, eventID uuid.UUID, isSettled bool) error {
	query := `UPDATE events SET is_settled = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(ctx, query, isSettled, eventID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("EVENT NOT FOUND")
	}
	return nil
}

func (h *eventRepository) Delete(ctx context.Context, eventID uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1;`
	result, err := h.db.Exec(ctx, query, eventID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("EVENT NOT FOUND")
	}
	return nil
}