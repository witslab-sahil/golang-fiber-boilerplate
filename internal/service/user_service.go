package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/witslab-sahil/fiber-boilerplate/internal/models"
	"github.com/witslab-sahil/fiber-boilerplate/internal/repository"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/logger"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword  = errors.New("invalid password")
)

type UserService interface {
	Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*models.UserResponse, int64, error)
	GetByID(ctx context.Context, id uint) (*models.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.UserResponse, error)
	Delete(ctx context.Context, id uint) error
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	repo            repository.UserRepository
	logger          logger.Logger
	tracer          trace.Tracer
	userCounter     metric.Int64Counter
	requestDuration metric.Float64Histogram
}

func NewUserService(repo repository.UserRepository, logger logger.Logger) UserService {
	meter := otel.Meter("user-service")
	
	userCounter, _ := meter.Int64Counter(
		"user_operations_total",
		metric.WithDescription("Total number of user operations"),
	)
	
	requestDuration, _ := meter.Float64Histogram(
		"user_operation_duration_seconds",
		metric.WithDescription("Duration of user operations in seconds"),
	)
	
	return &userService{
		repo:            repo,
		logger:          logger,
		tracer:          otel.Tracer("user-service"),
		userCounter:     userCounter,
		requestDuration: requestDuration,
	}
}

func (s *userService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	start := time.Now()
	ctx, span := s.tracer.Start(ctx, "UserService.Create")
	defer func() {
		span.End()
		duration := time.Since(start).Seconds()
		s.requestDuration.Record(ctx, duration, metric.WithAttributes(
			attribute.String("operation", "create"),
		))
	}()

	span.SetAttributes(
		attribute.String("user.email", req.Email),
		attribute.String("user.username", req.Username),
	)

	existingUser, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingUser != nil {
		span.SetAttributes(attribute.Bool("user.already_exists", true))
		return nil, ErrUserAlreadyExists
	}

	existingUser, err = s.repo.GetByUsername(req.Username)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to check existing username: %w", err)
	}
	if existingUser != nil {
		span.SetAttributes(attribute.Bool("user.already_exists", true))
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Roles:     req.Roles,
		IsActive:  true,
	}

	if err := s.repo.Create(user); err != nil {
		span.RecordError(err)
		s.logger.Errorf("Failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	span.SetAttributes(attribute.Int64("user.id", int64(user.ID)))
	s.userCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("operation", "create"),
		attribute.String("status", "success"),
	))
	s.logger.Infof("User created successfully: %s", user.Email)
	return user.ToResponse(), nil
}

func (s *userService) GetAll(ctx context.Context, page, pageSize int) ([]*models.UserResponse, int64, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.GetAll")
	defer span.End()

	span.SetAttributes(
		attribute.Int("pagination.page", page),
		attribute.Int("pagination.page_size", pageSize),
	)

	users, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	span.SetAttributes(
		attribute.Int64("users.total", total),
		attribute.Int("users.count", len(users)),
	)

	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.GetByID")
	defer span.End()

	span.SetAttributes(attribute.Int64("user.id", int64(id)))

	user, err := s.repo.GetByID(id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		span.SetAttributes(attribute.Bool("user.not_found", true))
		return nil, ErrUserNotFound
	}

	return user.ToResponse(), nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.GetByEmail")
	defer span.End()

	span.SetAttributes(attribute.String("user.email", email))

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		span.SetAttributes(attribute.Bool("user.not_found", true))
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Update")
	defer span.End()

	span.SetAttributes(attribute.Int64("user.id", int64(id)))
	user, err := s.repo.GetByID(id)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		span.SetAttributes(attribute.Bool("user.not_found", true))
		return nil, ErrUserNotFound
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.repo.GetByEmail(req.Email)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrUserAlreadyExists
		}
		user.Email = req.Email
	}

	if req.Username != "" && req.Username != user.Username {
		existingUser, err := s.repo.GetByUsername(req.Username)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to check existing username: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrUserAlreadyExists
		}
		user.Username = req.Username
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(user); err != nil {
		span.RecordError(err)
		s.logger.Errorf("Failed to update user: %v", err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Infof("User updated successfully: %s", user.Email)
	return user.ToResponse(), nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	ctx, span := s.tracer.Start(ctx, "UserService.Delete")
	defer span.End()

	span.SetAttributes(attribute.Int64("user.id", int64(id)))
	user, err := s.repo.GetByID(id)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		span.SetAttributes(attribute.Bool("user.not_found", true))
		return ErrUserNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		span.RecordError(err)
		s.logger.Errorf("Failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Infof("User deleted successfully: %s", user.Email)
	return nil
}

func (s *userService) CreateUser(user *models.User) (*models.User, error) {
	if err := s.repo.Create(user); err != nil {
		s.logger.Errorf("Failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}