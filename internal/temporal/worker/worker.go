package worker

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"github.com/sirupsen/logrus"

	"github.com/witslab-sahil/fiber-boilerplate/internal/temporal/activities"
	"github.com/witslab-sahil/fiber-boilerplate/internal/temporal/workflows"
)

type Worker struct {
	client client.Client
	worker worker.Worker
	logger *logrus.Logger
}

func NewWorker(c client.Client, taskQueue string, logger *logrus.Logger) (*Worker, error) {
	w := worker.New(c, taskQueue, worker.Options{
		MaxConcurrentActivityExecutionSize:     10,
		MaxConcurrentWorkflowTaskExecutionSize: 10,
	})

	// Register workflows
	w.RegisterWorkflow(workflows.UserOnboardingWorkflowFunc)

	// Register activities
	activityHandler := activities.NewActivities(logger)
	w.RegisterActivity(activityHandler.SendWelcomeEmail)
	w.RegisterActivity(activityHandler.SendFollowUpEmail)
	w.RegisterActivity(activityHandler.CreateUserProfile)
	w.RegisterActivity(activityHandler.SendPushNotification)
	w.RegisterActivity(activityHandler.SendSMSNotification)

	return &Worker{
		client: c,
		worker: w,
		logger: logger,
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	err := w.worker.Start()
	if err != nil {
		return fmt.Errorf("failed to start worker: %w", err)
	}

	w.logger.Info("Temporal worker started")

	// Wait for context cancellation
	<-ctx.Done()

	w.worker.Stop()
	w.logger.Info("Temporal worker stopped")

	return nil
}