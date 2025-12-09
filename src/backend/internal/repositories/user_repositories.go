package repositories

import (
	"context"
	"log"
	"sharebill-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (id, name, email, password_hash, bank_name, account_number, account_name) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	var insertedID uuid.UUID
	err := r.db.QueryRow(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.BankName, user.AccountNumber, user.AccountName).Scan(&insertedID)
	
	if(err != nil) {
		log.Printf("Error inserting user: %v", err)
		return err
	}

	log.Printf("Inserting successfully: %s", insertedID)
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, name, email, password_hash, bank_name, account_number, account_name, created_at, updated_at FROM users WHERE email=$1 LIMIT 1"
	
	row := r.db.QueryRow(ctx, query, email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.BankName, &user.AccountNumber, &user.AccountName, &user.CreatedAt, &user.UpdatedAt)
	if(err == pgx.ErrNoRows) {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	return user, nil
}