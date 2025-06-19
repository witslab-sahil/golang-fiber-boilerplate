package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/logger"
)

type OPAMiddleware struct {
	opaURL    string
	logger    logger.Logger
	jwtSecret string
}

type User struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type OPAInput struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	User   *User  `json:"user"`
}

type OPARequest struct {
	Input OPAInput `json:"input"`
}

type OPAResponse struct {
	Result bool `json:"result"`
}

func NewOPAMiddleware(opaURL string, logger logger.Logger) *OPAMiddleware {
	return &OPAMiddleware{
		opaURL: opaURL,
		logger: logger,
	}
}

func (m *OPAMiddleware) SetJWTSecret(secret string) {
	m.jwtSecret = secret
}

func (m *OPAMiddleware) Authorize() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Remove Bearer prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		// Parse and validate token
		user, err := m.parseToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Store user in context
		c.Locals("user", user)

		// Check authorization with OPA
		allowed, err := m.checkAuthorization(c.Method(), c.Path(), user)
		if err != nil {
			m.logger.Error("Failed to check authorization: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Authorization check failed",
			})
		}

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		return c.Next()
	}
}

func (m *OPAMiddleware) parseToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
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

		return &User{
			ID:    claims["sub"].(string),
			Email: claims["email"].(string),
			Roles: roleStrings,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (m *OPAMiddleware) checkAuthorization(method, path string, user *User) (bool, error) {
	// Create OPA request
	opaReq := OPARequest{
		Input: OPAInput{
			Method: method,
			Path:   path,
			User:   user,
		},
	}

	// Marshal request
	body, err := json.Marshal(opaReq)
	if err != nil {
		return false, err
	}

	// Send request to OPA
	resp, err := http.Post(
		fmt.Sprintf("%s/v1/data/authz/allow", m.opaURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse response
	var opaResp OPAResponse
	if err := json.NewDecoder(resp.Body).Decode(&opaResp); err != nil {
		return false, err
	}

	return opaResp.Result, nil
}