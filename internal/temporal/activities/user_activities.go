package activities

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"github.com/sirupsen/logrus"
)

type Activities struct {
	logger *logrus.Logger
}

func NewActivities(logger *logrus.Logger) *Activities {
	return &Activities{
		logger: logger,
	}
}

// Email activities
type SendEmailInput struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type SendEmailResult struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
}

func (a *Activities) SendWelcomeEmail(ctx context.Context, input SendEmailInput) (SendEmailResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending welcome email", "email", input.Email)

	// Simulate email sending
	// In production, integrate with email service (SendGrid, SES, etc.)
	time.Sleep(100 * time.Millisecond)

	return SendEmailResult{
		Success:   true,
		MessageID: fmt.Sprintf("welcome-%d-%d", input.UserID, time.Now().Unix()),
	}, nil
}

func (a *Activities) SendFollowUpEmail(ctx context.Context, input SendEmailInput) (SendEmailResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending follow-up email", "email", input.Email)

	// Simulate email sending
	time.Sleep(100 * time.Millisecond)

	return SendEmailResult{
		Success:   true,
		MessageID: fmt.Sprintf("followup-%d-%d", input.UserID, time.Now().Unix()),
	}, nil
}

// Profile activities
type CreateProfileInput struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type CreateProfileResult struct {
	Success   bool   `json:"success"`
	ProfileID string `json:"profile_id"`
}

func (a *Activities) CreateUserProfile(ctx context.Context, input CreateProfileInput) (CreateProfileResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Creating user profile", "userID", input.UserID)

	// Simulate profile creation
	// In production, this would interact with your database
	time.Sleep(200 * time.Millisecond)

	return CreateProfileResult{
		Success:   true,
		ProfileID: fmt.Sprintf("profile-%d", input.UserID),
	}, nil
}

// Notification activities
type NotificationInput struct {
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
}

type NotificationResult struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

func (a *Activities) SendPushNotification(ctx context.Context, input NotificationInput) (NotificationResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending push notification", "userID", input.UserID)

	// Simulate push notification
	// In production, integrate with FCM, APNS, etc.
	time.Sleep(150 * time.Millisecond)

	return NotificationResult{
		Success: true,
		ID:      fmt.Sprintf("push-%d-%d", input.UserID, time.Now().Unix()),
	}, nil
}

func (a *Activities) SendSMSNotification(ctx context.Context, input NotificationInput) (NotificationResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending SMS notification", "userID", input.UserID)

	// Simulate SMS sending
	// In production, integrate with Twilio, SNS, etc.
	time.Sleep(150 * time.Millisecond)

	return NotificationResult{
		Success: true,
		ID:      fmt.Sprintf("sms-%d-%d", input.UserID, time.Now().Unix()),
	}, nil
}

// Helper function to register all activities
func RegisterActivities(w interface {
	RegisterActivity(fn interface{}, options ...interface{})
}) {
	activities := &Activities{}
	w.RegisterActivity(activities.SendWelcomeEmail)
	w.RegisterActivity(activities.SendFollowUpEmail)
	w.RegisterActivity(activities.CreateUserProfile)
	w.RegisterActivity(activities.SendPushNotification)
	w.RegisterActivity(activities.SendSMSNotification)
}