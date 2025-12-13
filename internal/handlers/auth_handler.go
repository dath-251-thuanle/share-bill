package handlers

import (
	models "BACKEND/internal/dto"
	services "BACKEND/internal/services"
	"BACKEND/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *services.UserService
}
func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}


// POST /api/auth/register/request
// Body: { "email": "a@gmail.com" }
func (h *UserHandler) RegisterRequestCode(c *fiber.Ctx) error {
	var req models.RegisterRequestCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "INVALID_BODY", Message: "Invalid JSON format",
		})
	}

	err := h.service.RegisterRequestCode(c.Context(), req.Email)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Message: "Verification code sent to email",
	})
}

// POST /api/auth/register/confirm
// Body: { "email": "...", "code": "123456", "password": "...", "name": "..." }
func (h *UserHandler) RegisterConfirm(c *fiber.Ctx) error {
	var req models.RegisterConfirmRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "INVALID_BODY", Message: "Invalid JSON format",
		})
	}

	resp, err := h.service.RegisterConfirm(c.Context(), req)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    resp,
	})
}


// POST /api/auth/login
// Body: { "email": "...", "password": "..." }
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "INVALID_BODY", Message: "Invalid JSON format",
		})
	}

	resp, err := h.service.Login(c.Context(), req)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Message: "Login successfully",
		Data:    resp,
	})
}

// GET /api/v1/users/profile
// Header: Authorization: Bearer <token>
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Error: "UNAUTHORIZED", Message: "User ID not found in context",
		})
	}

	resp, err := h.service.GetProfile(c.Context(), userID)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Data:    resp,
	})
}

// PUT /api/v1/users/profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "INVALID_BODY", Message: "Invalid JSON format",
		})
	}

	resp, err := h.service.UpdateProfile(c.Context(), userID, req)
	if err != nil {
		return utils.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    resp,
	})
}

func (h *UserHandler) UpdateAvatar(c *fiber.Ctx) error{
	userID := c.Locals("user_id").(int64)
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "INVALID_BODY",
			Message: "Image is required",
		})
	}
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "OPEN_FILE_FAILED",
			Message:   "Unable to open file stream",
		})	
	}
	defer file.Close()

	resp, err := h.service.UpdateUserAvatar(c.UserContext(), userID, file, fileHeader.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: "UPLOAD_FILE_FAILED",
			Message:   "Upload failed: " + err.Error(),
		})	
	}
	return c.Status(fiber.StatusOK).JSON(models.SuccessResponse{
		Success: true,
		Data: resp,
	})
}
