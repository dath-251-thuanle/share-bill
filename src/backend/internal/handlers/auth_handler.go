package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"sharebill-backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Name          string  `json:"name"`
		Email         string  `json:"email"`
		Password      string  `json:"password"`
		BankName      *string `json:"bank_name"`
		AccountNumber *string `json:"account_number"`
		AccountName   *string `json:"account_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name, email, and password are required"})
	}

	user, token, err := h.authService.Register(c.Context(), req.Name, req.Email, req.Password, req.BankName, req.AccountNumber, req.AccountName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	user, token, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	return c.JSON(fiber.Map{
		"message": "Login successfully",
		"token": token,
		"user": user,
	})
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userIDstr, ok := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDstr)

	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	user, err := h.authService.GetUserByID(c.Context(), userID)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"message": "Profile fetched successfully",
		"profile": fiber.Map{
			"name":           user.Name,
			"email":          user.Email,
			"bank_name":      user.BankName,
			"account_number": user.AccountNumber,
			"account_name":   user.AccountName,
		},
	})
}

func (h *AuthHandler) EditProfile(c *fiber.Ctx) error {
	userIDstr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	var req struct {
		Name          string `json:"name"`
		BankName      *string `json:"bank_name"`
		AccountNumber *string `json:"account_number"`
		AccountName   *string `json:"account_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	updatedUser, str, err := h.authService.EditProfile(c.Context(), userID, req.Name, req.BankName, req.AccountNumber, req.AccountName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    updatedUser,
		"token":   str,
	})
}

func (h *AuthHandler) DeleteAccount(c *fiber.Ctx) error {
	userIDstr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	var req struct {
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Wrong password"})
	}
	if req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}

	err = h.authService.DeleteUser(c.Context(), userID, req.Password)
	if err != nil {
		if err.Error() == "incorrect password" {
			return c.Status(400).JSON(fiber.Map{"error": "Wrong password"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}
