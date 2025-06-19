package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/witslab-sahil/fiber-boilerplate/internal/config"
	"github.com/witslab-sahil/fiber-boilerplate/internal/temporal/worker"
	"github.com/witslab-sahil/fiber-boilerplate/internal/temporal/workflows"
	pkgTemporal "github.com/witslab-sahil/fiber-boilerplate/pkg/temporal"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal("Invalid log level:", err)
	}
	logger.SetLevel(level)

	// Create Temporal client
	temporalHost := os.Getenv("TEMPORAL_HOST")
	if temporalHost == "" {
		temporalHost = "localhost:7233"
	}

	temporalClient, err := pkgTemporal.NewClient(temporalHost, cfg.TemporalNamespace, cfg.OtelEnabled)
	if err != nil {
		logger.Fatal("Failed to create Temporal client:", err)
	}
	defer temporalClient.Close()

	// Create worker
	taskQueue := os.Getenv("TASK_QUEUE")
	if taskQueue == "" {
		taskQueue = workflows.OnboardingTaskQueue
	}

	w, err := worker.NewWorker(temporalClient.GetClient(), taskQueue, logger)
	if err != nil {
		logger.Fatal("Failed to create worker:", err)
	}

	// Start worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := w.Start(ctx); err != nil {
			logger.Fatal("Failed to start worker:", err)
		}
	}()

	logger.Infof("Worker started for task queue: %s", taskQueue)

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("Shutting down worker...")
	cancel()
}