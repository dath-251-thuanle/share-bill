package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name	  string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	OwnerID   uuid.UUID `json:"owner_id" db:"owner_id"`
	IsSettled  bool      `json:"is_settled" db:"is_settled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateEventRequest struct {
	Name 	  string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"max=500"`
}

type UpdateEventRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	IsSettled   *bool   `json:"is_settled,omitempty" validate:"omitempty"`
}

type EventDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	OwnerID     uuid.UUID `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	OwnerEmail  string    `json:"owner_email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	OwnerID     uuid.UUID `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`     // tên chủ event (join từ users)
	IsSettled   bool      `json:"is_settled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TotalExpense float64  `json:"total_expense"` // tổng chi phí (tính từ bảng expenses)
	MemberCount  int      `json:"member_count"`  // số người tham gia
}