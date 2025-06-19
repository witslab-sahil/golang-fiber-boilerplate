package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/witslab-sahil/fiber-boilerplate/internal/config"
	"github.com/witslab-sahil/fiber-boilerplate/internal/handlers"
	"github.com/witslab-sahil/fiber-boilerplate/internal/middleware"
	"github.com/witslab-sahil/fiber-boilerplate/internal/models"
	opaMiddleware "github.com/witslab-sahil/fiber-boilerplate/internal/opa/middleware"
	"github.com/witslab-sahil/fiber-boilerplate/internal/repository"
	"github.com/witslab-sahil/fiber-boilerplate/internal/service"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/database"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/logger"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/telemetry"
	pkgTemporal "github.com/witslab-sahil/fiber-boilerplate/pkg/temporal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()
	
	logger := logger.New(cfg.LogLevel)
	logger.Info("Starting application...")

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}

	// Run migrations
	if err := database.Migrate(db, &models.User{}); err != nil {
		logger.Fatal("Failed to run migrations: ", err)
	}

	// Initialize telemetry
	var tel *telemetry.Telemetry
	if cfg.OtelEnabled {
		telCfg := &telemetry.Config{
			ServiceName:    cfg.OtelServiceName,
			ServiceVersion: "1.0.0",
			Environment:    cfg.Environment,
			ExporterType:   cfg.OtelExporterType,
			Endpoint:       cfg.OtelEndpoint,
			Enabled:        cfg.OtelEnabled,
		}
		
		var err error
		tel, err = telemetry.New(context.Background(), telCfg)
		if err != nil {
			logger.Fatal("Failed to initialize telemetry: ", err)
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := tel.Shutdown(ctx); err != nil {
				logger.Error("Error shutting down telemetry: ", err)
			}
		}()
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, logger)

	// Initialize Temporal client
	var temporalClient *pkgTemporal.Client
	if cfg.TemporalHost != "" {
		temporalClient, err = pkgTemporal.NewClient(cfg.TemporalHost, cfg.TemporalNamespace, cfg.OtelEnabled)
		if err != nil {
			logger.Warn("Failed to connect to Temporal, workflows will be disabled: ", err)
		} else {
			defer temporalClient.Close()
		}
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: true,
	}))
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(logger))

	// Initialize telemetry middleware
	if cfg.OtelEnabled {
		app.Use(middleware.Tracing(cfg.OtelServiceName))
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, logger, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userService, logger)
	workflowHandler := handlers.NewWorkflowHandler(temporalClient, logger)

	// Health check
	app.Get("/health", handlers.HealthCheck)

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	if cfg.OPAEnabled {
		// Initialize OPA middleware
		opaMiddleware := opaMiddleware.NewOPAMiddleware(cfg.OPAURL, logger)
		api.Use(opaMiddleware.Authorize())
	}

	// User routes (protected)
	users := api.Group("/users")
	users.Get("/", userHandler.GetAll)
	users.Get("/:id", userHandler.GetByID)
	users.Post("/", userHandler.Create)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	// Workflow routes (protected)
	if temporalClient != nil {
		workflows := api.Group("/workflows")
		workflows.Post("/user-onboarding", workflowHandler.StartUserOnboarding)
		workflows.Get("/user-onboarding/:id/status", workflowHandler.GetWorkflowStatus)
		workflows.Get("/", workflowHandler.ListWorkflows)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			logger.Fatal("Failed to start server: ", err)
		}
	}()

	logger.Info("Server started on port ", cfg.Port)

	<-quit
	logger.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exited")
}