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
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
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

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := "SELECT id, name, email, password_hash, bank_name, account_number, account_name, created_at, updated_at FROM users WHERE id=$1 LIMIT 1"
	
	row := r.db.QueryRow(ctx, query, id)
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

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name = $1, bank_name = $2, account_number = $3, account_name = $4, updated_at = NOW() WHERE id = $5`
	_, err := r.db.Exec(ctx, query, user.Name, user.BankName, user.AccountNumber, user.AccountName, user.ID,)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id=$1"
	_, err := r.db.Exec(ctx, query, id)
	return err
}