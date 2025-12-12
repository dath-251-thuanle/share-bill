package models
import (
	"time"
	"github.com/google/uuid"
)

type Participant struct {
	UserID    uuid.UUID      `db:"event_id"`
	EventID   uint      `db:"user_id"`
	JoinedAt  time.Time `db:"joined_at"`
}

type ParticipantResponse struct {
	UserID   uuid.UUID   `json:"user_id"`
	Name	 string `json:"name"`
	Email    string `json:"email"`
	BankName string `json:"bank_name,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	AccountName   string `json:"account_name,omitempty"`
	JoinedAt time.Time `json:"joined_at"`
}