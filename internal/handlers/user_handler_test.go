package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/witslab-sahil/fiber-boilerplate/internal/models"
	"github.com/witslab-sahil/fiber-boilerplate/internal/service"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/logger"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) GetAll(ctx context.Context, page, pageSize int) ([]*models.UserResponse, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.UserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) GetByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...)...)
}

func (m *MockLogger) WithFields(fields map[string]interface{}) logger.Logger {
	m.Called(fields)
	return m
}

func TestUserHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		mockLogger := new(MockLogger)
		handler := NewUserHandler(mockService, mockLogger)
		app := fiber.New()
		
		req := &models.CreateUserRequest{
			Email:     "test@example.com",
			Username:  "testuser",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		expectedUser := &models.UserResponse{
			ID:        1,
			Email:     req.Email,
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		mockService.On("Create", mock.Anything, req).Return(expectedUser, nil)

		app.Post("/users", handler.Create)

		reqBody, _ := json.Marshal(req)
		request := httptest.NewRequest("POST", "/users", bytes.NewReader(reqBody))
		request.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(request)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

		var response models.UserResponse
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &response)

		assert.Equal(t, expectedUser.ID, response.ID)
		assert.Equal(t, expectedUser.Email, response.Email)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockService := new(MockUserService)
		mockLogger := new(MockLogger)
		handler := NewUserHandler(mockService, mockLogger)
		app := fiber.New()
		app.Post("/users", handler.Create)

		request := httptest.NewRequest("POST", "/users", bytes.NewReader([]byte("invalid json")))
		request.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(request)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		mockLogger := new(MockLogger)
		handler := NewUserHandler(mockService, mockLogger)
		app := fiber.New()
		
		req := &models.CreateUserRequest{
			Email:    "test@example.com",
			Username: "testuser",
		}

		mockService.On("Create", mock.Anything, req).Return(nil, errors.New("service error"))
		mockLogger.On("Error", "Failed to create user: ", mock.Anything).Return()

		app.Post("/users", handler.Create)

		reqBody, _ := json.Marshal(req)
		request := httptest.NewRequest("POST", "/users", bytes.NewReader(reqBody))
		request.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(request)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		mockLogger := new(MockLogger)
		handler := NewUserHandler(mockService, mockLogger)
		app := fiber.New()
		
		expectedUser := &models.UserResponse{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
		}

		mockService.On("GetByID", mock.Anything, uint(1)).Return(expectedUser, nil)

		app.Get("/users/:id", handler.GetByID)

		request := httptest.NewRequest("GET", "/users/1", nil)
		resp, _ := app.Test(request)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var response models.UserResponse
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &response)

		assert.Equal(t, expectedUser.ID, response.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockService := new(MockUserService)
		mockLogger := new(MockLogger)
		handler := NewUserHandler(mockService, mockLogger)
		app := fiber.New()
		
		mockService.On("GetByID", mock.Anything, uint(999)).Return(nil, service.ErrUserNotFound)

		app.Get("/users/:id", handler.GetByID)

		request := httptest.NewRequest("GET", "/users/999", nil)
		resp, _ := app.Test(request)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}