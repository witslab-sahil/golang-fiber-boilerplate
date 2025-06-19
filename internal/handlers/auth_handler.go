package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/witslab-sahil/fiber-boilerplate/internal/models"
	"github.com/witslab-sahil/fiber-boilerplate/internal/opa/middleware"
	"github.com/witslab-sahil/fiber-boilerplate/internal/service"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/logger"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/utils"
)

type AuthHandler struct {
	userService service.UserService
	logger      logger.Logger
	jwtSecret   string
}

func NewAuthHandler(userService service.UserService, logger logger.Logger, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		logger:      logger,
		jwtSecret:   jwtSecret,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create user
	user, err := h.userService.Create(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create user: ", err)
		if err.Error() == "email already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		h.logger.Error("Failed to generate token: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user by email
	user, err := h.userService.GetByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT token
	token, err := h.generateToken(&models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Roles:     user.Roles,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	if err != nil {
		h.logger.Error("Failed to generate token: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"user": models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     user.Roles,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		"token": token,
	})
}

func (h *AuthHandler) generateToken(user *models.UserResponse) (string, error) {
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", user.ID),
		"email": user.Email,
		"roles": user.Roles,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// ParseToken parses and validates JWT token
func (h *AuthHandler) ParseToken(tokenString string) (*middleware.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		roles, _ := claims["roles"].([]interface{})
		roleStrings := make([]string, len(roles))
		for i, r := range roles {
			roleStrings[i] = r.(string)
		}

		return &middleware.User{
			ID:    claims["sub"].(string),
			Email: claims["email"].(string),
			Roles: roleStrings,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}