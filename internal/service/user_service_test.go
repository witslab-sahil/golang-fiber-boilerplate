package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/witslab-sahil/fiber-boilerplate/internal/models"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/logger"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(page, pageSize int) ([]*models.User, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]*models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{}) {}
func (m *MockLogger) Debugf(format string, args ...interface{}) {}
func (m *MockLogger) Info(args ...interface{}) {}
func (m *MockLogger) Infof(format string, args ...interface{}) {}
func (m *MockLogger) Warn(args ...interface{}) {}
func (m *MockLogger) Warnf(format string, args ...interface{}) {}
func (m *MockLogger) Error(args ...interface{}) {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}
func (m *MockLogger) Fatal(args ...interface{}) {}
func (m *MockLogger) Fatalf(format string, args ...interface{}) {}
func (m *MockLogger) WithFields(fields map[string]interface{}) logger.Logger {
	return m
}

func TestUserService_Create(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	t.Run("Success", func(t *testing.T) {
		req := &models.CreateUserRequest{
			Email:     "test@example.com",
			Username:  "testuser",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		mockRepo.On("GetByEmail", req.Email).Return(nil, nil).Once()
		mockRepo.On("GetByUsername", req.Username).Return(nil, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()

		result, err := service.Create(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		req := &models.CreateUserRequest{
			Email:    "existing@example.com",
			Username: "newuser",
			Password: "password123",
		}

		existingUser := &models.User{Email: req.Email}
		mockRepo.On("GetByEmail", req.Email).Return(existingUser, nil).Once()

		result, err := service.Create(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, ErrUserAlreadyExists, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Username Already Exists", func(t *testing.T) {
		req := &models.CreateUserRequest{
			Email:    "new@example.com",
			Username: "existinguser",
			Password: "password123",
		}

		existingUser := &models.User{Username: req.Username}
		mockRepo.On("GetByEmail", req.Email).Return(nil, nil).Once()
		mockRepo.On("GetByUsername", req.Username).Return(existingUser, nil).Once()

		result, err := service.Create(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, ErrUserAlreadyExists, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
		}

		mockRepo.On("GetByID", uint(1)).Return(user, nil).Once()

		result, err := service.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Email, result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, nil).Once()

		result, err := service.GetByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
		}

		req := &models.UpdateUserRequest{
			FirstName: "Updated",
			LastName:  "Name",
		}

		mockRepo.On("GetByID", uint(1)).Return(user, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil).Once()

		result, err := service.Update(context.Background(), 1, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.FirstName, result.FirstName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		req := &models.UpdateUserRequest{
			FirstName: "Updated",
		}

		mockRepo.On("GetByID", uint(999)).Return(nil, nil).Once()

		result, err := service.Update(context.Background(), 999, req)
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		user := &models.User{
			ID:    1,
			Email: "current@example.com",
		}

		req := &models.UpdateUserRequest{
			Email: "existing@example.com",
		}

		existingUser := &models.User{
			ID:    2,
			Email: req.Email,
		}

		mockRepo.On("GetByID", uint(1)).Return(user, nil).Once()
		mockRepo.On("GetByEmail", req.Email).Return(existingUser, nil).Once()

		result, err := service.Update(context.Background(), 1, req)
		assert.Error(t, err)
		assert.Equal(t, ErrUserAlreadyExists, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			ID:    1,
			Email: "test@example.com",
		}

		mockRepo.On("GetByID", uint(1)).Return(user, nil).Once()
		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		err := service.Delete(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return(nil, nil).Once()

		err := service.Delete(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetAll(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	t.Run("Success", func(t *testing.T) {
		users := []*models.User{
			{ID: 1, Email: "test1@example.com"},
			{ID: 2, Email: "test2@example.com"},
		}

		mockRepo.On("GetAll", 1, 10).Return(users, int64(2), nil).Once()

		results, total, err := service.GetAll(context.Background(), 1, 10)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("GetAll", 1, 10).Return([]*models.User{}, int64(0), errors.New("database error")).Once()

		results, total, err := service.GetAll(context.Background(), 1, 10)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Equal(t, int64(0), total)
		mockRepo.AssertExpectations(t)
	})
}