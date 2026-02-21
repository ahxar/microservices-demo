package service

import (
	"context"
	"fmt"
	"log"

	"github.com/safar/microservices-demo/services/notification/internal/repository"
	"github.com/safar/microservices-demo/services/notification/internal/smtp"
	"github.com/safar/microservices-demo/services/notification/internal/templates"
)

type NotificationService struct {
	repo       *repository.NotificationRepository
	smtpClient *smtp.SMTPClient
}

func NewNotificationService(repo *repository.NotificationRepository, smtpClient *smtp.SMTPClient) *NotificationService {
	return &NotificationService{
		repo:       repo,
		smtpClient: smtpClient,
	}
}

func (s *NotificationService) SendOrderConfirmation(ctx context.Context, data templates.OrderConfirmationData, email string) error {
	// Render email template
	body, err := templates.RenderOrderConfirmation(data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	subject := fmt.Sprintf("Order Confirmation - #%s", data.OrderID)

	// Create notification record
	notification := &repository.Notification{
		RecipientEmail: email,
		RecipientName:  data.FirstName,
		Type:           "order_confirmation",
		Subject:        subject,
		Body:           body,
		Status:         "pending",
		RetryCount:     0,
	}

	createdNotification, err := s.repo.CreateNotification(notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Send email asynchronously
	go s.sendEmailAsync(createdNotification.ID, email, subject, body)

	return nil
}

func (s *NotificationService) SendShippingUpdate(ctx context.Context, data templates.ShippingUpdateData, email string) error {
	// Render email template
	body, err := templates.RenderShippingUpdate(data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	subject := fmt.Sprintf("Shipping Update - Order #%s", data.OrderID)

	// Create notification record
	notification := &repository.Notification{
		RecipientEmail: email,
		RecipientName:  data.FirstName,
		Type:           "shipping_update",
		Subject:        subject,
		Body:           body,
		Status:         "pending",
		RetryCount:     0,
	}

	createdNotification, err := s.repo.CreateNotification(notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Send email asynchronously
	go s.sendEmailAsync(createdNotification.ID, email, subject, body)

	return nil
}

func (s *NotificationService) SendWelcomeEmail(ctx context.Context, data templates.WelcomeData, email string) error {
	// Render email template
	body, err := templates.RenderWelcome(data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	subject := "Welcome to Our Platform!"

	// Create notification record
	notification := &repository.Notification{
		RecipientEmail: email,
		RecipientName:  data.FirstName,
		Type:           "welcome",
		Subject:        subject,
		Body:           body,
		Status:         "pending",
		RetryCount:     0,
	}

	createdNotification, err := s.repo.CreateNotification(notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Send email asynchronously
	go s.sendEmailAsync(createdNotification.ID, email, subject, body)

	return nil
}

func (s *NotificationService) SendPasswordReset(ctx context.Context, data templates.PasswordResetData, email string) error {
	// Render email template
	body, err := templates.RenderPasswordReset(data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	subject := "Password Reset Request"

	// Create notification record
	notification := &repository.Notification{
		RecipientEmail: email,
		RecipientName:  data.FirstName,
		Type:           "password_reset",
		Subject:        subject,
		Body:           body,
		Status:         "pending",
		RetryCount:     0,
	}

	createdNotification, err := s.repo.CreateNotification(notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Send email asynchronously
	go s.sendEmailAsync(createdNotification.ID, email, subject, body)

	return nil
}

func (s *NotificationService) sendEmailAsync(notificationID, to, subject, body string) {
	err := s.smtpClient.SendEmail(to, subject, body)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		errorMsg := err.Error()
		s.repo.UpdateStatus(notificationID, "failed", &errorMsg)
		s.repo.IncrementRetryCount(notificationID)
		return
	}

	log.Printf("Email sent successfully to %s", to)
	s.repo.UpdateStatus(notificationID, "sent", nil)
}

func (s *NotificationService) GetNotification(ctx context.Context, id string) (*repository.Notification, error) {
	return s.repo.GetNotification(id)
}

func (s *NotificationService) ListNotifications(ctx context.Context, email string, page, pageSize int) ([]*repository.Notification, int, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListNotifications(email, pageSize, offset)
}
