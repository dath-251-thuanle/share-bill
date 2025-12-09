package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID		uuid.UUID `json:"id" db:"id"`
	Name 	string    `json:"name" db:"name"`
	Email 	string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	BankName	*string    `json:"bank_name" db:"bank_name"`
	AccountNumber *string    `json:"account_number" db:"account_number"`
	AccountName   *string    `json:"account_name" db:"account_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}