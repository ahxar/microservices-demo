package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type NotificationRepository struct {
	db *sql.DB
}

type Notification struct {
	ID             string
	RecipientEmail string
	RecipientName  string
	Type           string
	Subject        string
	Body           string
	Status         string
	RetryCount     int
	ErrorMessage   sql.NullString
	SentAt         sql.NullTime
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewNotificationRepository(databaseURL string) (*NotificationRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &NotificationRepository{db: db}, nil
}

func (r *NotificationRepository) Close() error {
	return r.db.Close()
}

func (r *NotificationRepository) CreateNotification(notification *Notification) (*Notification, error) {
	query := `
		INSERT INTO notifications (
			recipient_email, recipient_name, notification_type, subject, body, status, retry_count
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		notification.RecipientEmail,
		notification.RecipientName,
		notification.Type,
		notification.Subject,
		notification.Body,
		notification.Status,
		notification.RetryCount,
	).Scan(&notification.ID, &notification.CreatedAt, &notification.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

func (r *NotificationRepository) UpdateStatus(id, status string, errorMessage *string) error {
	var query string
	var err error

	if errorMessage != nil {
		query = `
			UPDATE notifications
			SET status = $1, error_message = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $3
		`
		_, err = r.db.Exec(query, status, *errorMessage, id)
	} else {
		query = `
			UPDATE notifications
			SET status = $1, sent_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`
		_, err = r.db.Exec(query, status, id)
	}

	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	return nil
}

func (r *NotificationRepository) IncrementRetryCount(id string) error {
	query := `UPDATE notifications SET retry_count = retry_count + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetNotification(id string) (*Notification, error) {
	notification := &Notification{}
	query := `
		SELECT id, recipient_email, recipient_name, notification_type, subject, body,
		       status, retry_count, error_message, sent_at, created_at, updated_at
		FROM notifications
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.RecipientEmail,
		&notification.RecipientName,
		&notification.Type,
		&notification.Subject,
		&notification.Body,
		&notification.Status,
		&notification.RetryCount,
		&notification.ErrorMessage,
		&notification.SentAt,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return notification, nil
}

func (r *NotificationRepository) ListNotifications(email string, limit, offset int) ([]*Notification, int, error) {
	var query string
	var args []interface{}

	if email != "" {
		query = `
			SELECT id, recipient_email, recipient_name, notification_type, subject, body,
			       status, retry_count, error_message, sent_at, created_at, updated_at
			FROM notifications
			WHERE recipient_email = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{email, limit, offset}
	} else {
		query = `
			SELECT id, recipient_email, recipient_name, notification_type, subject, body,
			       status, retry_count, error_message, sent_at, created_at, updated_at
			FROM notifications
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		notification := &Notification{}
		err := rows.Scan(
			&notification.ID,
			&notification.RecipientEmail,
			&notification.RecipientName,
			&notification.Type,
			&notification.Subject,
			&notification.Body,
			&notification.Status,
			&notification.RetryCount,
			&notification.ErrorMessage,
			&notification.SentAt,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	// Get total count
	var total int
	if email != "" {
		r.db.QueryRow("SELECT COUNT(*) FROM notifications WHERE recipient_email = $1", email).Scan(&total)
	} else {
		r.db.QueryRow("SELECT COUNT(*) FROM notifications").Scan(&total)
	}

	return notifications, total, nil
}
