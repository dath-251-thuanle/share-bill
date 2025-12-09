package services

import (
	"context"
	"errors"
	"os"
	"time"

	"sharebill-backend/internal/models"
	"sharebill-backend/internal/repositories"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string, bankName, accountNumber, accountName *string) (*models.User, string, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	
	if err != nil {
		return nil, "", err
	}

	if existingUser != nil {
		return nil, "", errors.New("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	newUser := &models.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		BankName:    bankName,
		AccountNumber: accountNumber,
		AccountName:   accountName,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": newUser.ID.String(),
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", err
	}

	return newUser, tokenString, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User ,string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	//check email
	if(user == nil || err != nil) {
		return nil, "", errors.New("email or password is wrong")
	}
	
	//check json
	if(user.Email == "" || user.PasswordHash == ""){
		return nil, "", errors.New("email or password must not be empty")
	}
	
	//check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("email or password is wrong")
	}

	// 3. Tạo JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"name":    user.Name,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}