package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"github.com/witslab-sahil/fiber-boilerplate/internal/temporal/activities"
)

const (
	UserOnboardingWorkflow = "UserOnboardingWorkflow"
	OnboardingTaskQueue    = "user-onboarding"
)

type UserOnboardingInput struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserOnboardingResult struct {
	Success           bool   `json:"success"`
	WelcomeEmailSent  bool   `json:"welcome_email_sent"`
	ProfileCreated    bool   `json:"profile_created"`
	NotificationsSent bool   `json:"notifications_sent"`
	Message           string `json:"message"`
}

func UserOnboardingWorkflowFunc(ctx workflow.Context, input UserOnboardingInput) (UserOnboardingResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting user onboarding workflow", "userID", input.UserID)

	result := UserOnboardingResult{
		Success: true,
	}

	// Activity options with retries
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    100 * time.Second,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Send welcome email
	var emailResult activities.SendEmailResult
	activityHandler := &activities.Activities{}
	err := workflow.ExecuteActivity(ctx, activityHandler.SendWelcomeEmail, activities.SendEmailInput{
		UserID: input.UserID,
		Email:  input.Email,
		Name:   input.Username,
	}).Get(ctx, &emailResult)
	if err != nil {
		logger.Error("Failed to send welcome email", "error", err)
		result.Success = false
		result.Message = "Failed to send welcome email"
	} else {
		result.WelcomeEmailSent = true
	}

	// Step 2: Create user profile
	var profileResult activities.CreateProfileResult
	err = workflow.ExecuteActivity(ctx, activityHandler.CreateUserProfile, activities.CreateProfileInput{
		UserID:   input.UserID,
		Username: input.Username,
	}).Get(ctx, &profileResult)
	if err != nil {
		logger.Error("Failed to create user profile", "error", err)
		result.Success = false
		result.Message += "; Failed to create user profile"
	} else {
		result.ProfileCreated = true
	}

	// Step 3: Send notifications (parallel execution)
	var futures []workflow.Future

	// Send push notification
	pushFuture := workflow.ExecuteActivity(ctx, activityHandler.SendPushNotification, activities.NotificationInput{
		UserID:  input.UserID,
		Message: "Welcome to our platform!",
	})
	futures = append(futures, pushFuture)

	// Send SMS notification
	smsFuture := workflow.ExecuteActivity(ctx, activityHandler.SendSMSNotification, activities.NotificationInput{
		UserID:  input.UserID,
		Message: "Welcome! Your account is ready.",
	})
	futures = append(futures, smsFuture)

	// Wait for all notifications to complete
	allNotificationsSent := true
	for _, future := range futures {
		var notificationResult activities.NotificationResult
		if err := future.Get(ctx, &notificationResult); err != nil {
			logger.Error("Failed to send notification", "error", err)
			allNotificationsSent = false
		}
	}
	result.NotificationsSent = allNotificationsSent

	// Step 4: Schedule follow-up (using timer)
	err = workflow.Sleep(ctx, 24*time.Hour) // Wait 24 hours
	if err != nil {
		logger.Error("Timer failed", "error", err)
		return result, err
	}

	// Send follow-up email
	err = workflow.ExecuteActivity(ctx, activityHandler.SendFollowUpEmail, activities.SendEmailInput{
		UserID: input.UserID,
		Email:  input.Email,
		Name:   input.Username,
	}).Get(ctx, &emailResult)
	if err != nil {
		logger.Error("Failed to send follow-up email", "error", err)
	}

	if result.Success {
		result.Message = "User onboarding completed successfully"
	}

	return result, nil
}