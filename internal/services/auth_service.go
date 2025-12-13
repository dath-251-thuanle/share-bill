package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"BACKEND/internal/config"
	"BACKEND/internal/db"
	database "BACKEND/internal/db/sqlc"
	models "BACKEND/internal/dto"
	"BACKEND/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	store       database.Store
	tokenMaker  *utils.JWTMaker
	config      config.Config
	redisClient *db.RedisClient 
	emailSender utils.EmailSender 
	uploadService *UploadService 
}
func NewUserService(
	store database.Store,
	tokenMaker *utils.JWTMaker,
	cfg config.Config,
	redisClient *db.RedisClient,
	emailSender utils.EmailSender,
	uploadService *UploadService,
) *UserService {
	return &UserService{
		store:       store,
		tokenMaker:  tokenMaker,
		config:      cfg,
		redisClient: redisClient,
		emailSender: emailSender,
		uploadService: uploadService,
	}
}

func (s *UserService) RegisterRequestCode(ctx context.Context, email string) error {
	_, err := s.store.GetUserByEmail(ctx, &email)
	if err == nil {
		return utils.ErrAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return utils.ErrInternalDB
	}

	otp := utils.GenerateOTP()

	// Key: otp:a@gmail.com  Value: 123456
	err = s.redisClient.SetOTP(ctx, email, otp, 5*time.Minute)
	if err != nil {
		return errors.New("failed to store otp in cache")
	}

	subject := "Sharever - Mã xác thực đăng ký"
	content := fmt.Sprintf(`
		<h1>Xin chào!</h1>
		<p>Mã xác thực của bạn là: <strong>%s</strong></p>
		<p>Mã này sẽ hết hạn sau 5 phút.</p>
	`, otp)
	to := []string{email}
	err = s.emailSender.SendEmail(subject, content, to)
	if err != nil {
		return errors.New("failed to send email")
	}

	return nil
}

func (s *UserService) RegisterConfirm(ctx context.Context, req models.RegisterConfirmRequest) (models.UserResponse, error) {
	storedOTP, err := s.redisClient.GetOTP(ctx, req.Email)
	if err == redis.Nil {
		return models.UserResponse{}, errors.New("otp expired or not found")
	} else if err != nil {
		return models.UserResponse{}, utils.ErrNotFound
	}
	if storedOTP != req.Code {
		return models.UserResponse{}, errors.New("invalid verification code")
	}
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.UserResponse{}, utils.ErrInternalDB
	}

	arg := database.CreateUserParams{
		Name:        req.Name,
		Email:       &req.Email,
		Password:    hashedPwd,
		PhoneNumber: nil,
		BankName:    nil,
		BankAccount: nil,
		BankOwner:   nil,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		return models.UserResponse{}, utils.ErrInternalDB
	}

	_ = s.redisClient.DeleteOTP(ctx, req.Email)
	//return s.mapUserResponse(user), nil
	return models.UserResponse{
		ID:          user.UserUuid.String(),
		Name:        user.Name,
		Email:       *user.Email,
		PhoneNumber: nil,
		Avatar:      nil, 
		BankInfo:    nil,           
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

// Login
func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	user, err := s.store.GetUserByEmail(ctx, &req.Email)
	if err != nil {
		return models.LoginResponse{}, utils.ErrNotFound
	}
	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return models.LoginResponse{}, utils.ErrUnauthorized // Sai pass
	}

	token, err := s.tokenMaker.CreateToken(user.UserID, s.config.TokenDuration)
	if err != nil {
		return models.LoginResponse{}, utils.ErrInternalDB
	}

	return models.LoginResponse{
		Token: token,
		User:  s.mapUserResponse(user),
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (models.UserResponse, error) {
	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return models.UserResponse{}, utils.ErrNotFound
	}
	return s.mapUserResponse(user), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int64, req models.UpdateProfileRequest) (models.UserResponse, error) {
	arg := database.UpdateUserParams{
		UserID: userID,
		Name:        s.toPtr(req.Name),
		Avatar: 	s.toPtr(req.Avatar),
		PhoneNumber: s.toPtr(req.PhoneNumber),
		BankName:    s.toPtr(req.BankName),
		BankAccount: s.toPtr(req.AccountNumber),
		BankOwner:   s.toPtr(req.AccountName),
	}

	user, err := s.store.UpdateUser(ctx, arg)
	if err != nil {
		return models.UserResponse{}, utils.ErrInternalDB
	}
	return s.mapUserResponse(user), nil
}

func (s *UserService) mapUserResponse(user database.User) models.UserResponse {
	var bankInfo *models.BankInfoDTO

	if user.BankName != nil {
		bankInfo = &models.BankInfoDTO{
			BankName:      *user.BankName,
			AccountNumber: *user.BankAccount,
			AccountName:   *user.BankOwner,
		}
	}
	
	return models.UserResponse{
		ID:          user.UserUuid.String(),
		Name:        user.Name,
		Email:       *user.Email,
		PhoneNumber: user.PhoneNumber,
		Avatar:      user.Avatar, 
		BankInfo:    bankInfo,           
		UpdatedAt:   user.UpdatedAt,
	}
}

func (s *UserService) toPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func (s *UserService) UpdateUserAvatar(ctx context.Context, userID int64, file multipart.File, filename string) (models.UserResponse, error) {
	avatarUrl, err := s.uploadService.UploadImage(ctx, file, filename)
	if err != nil {
		return models.UserResponse{}, err
	}

	updatedUser, err := s.store.UpdateUserAvatar(ctx, database.UpdateUserAvatarParams{
		UserID:     userID,
		Avatar: utils.StringToPtr(avatarUrl), 
	})
	if err != nil {
		return models.UserResponse{}, utils.ErrInternalDB
	}
	return models.UserResponse{
		ID:          updatedUser.UserUuid.String(),
		Name:        updatedUser.Name,
		Email:       *updatedUser.Email,
		PhoneNumber: updatedUser.PhoneNumber,
		Avatar:      updatedUser.Avatar, 
		BankInfo:    &models.BankInfoDTO{
			BankName: *updatedUser.BankName,
			AccountNumber: *updatedUser.BankAccount,
			AccountName: *updatedUser.BankOwner,
		},          
		UpdatedAt:   updatedUser.UpdatedAt,
	}, nil
}