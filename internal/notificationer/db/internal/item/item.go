package item

import (
	"github.com/google/uuid"
	"notification-scheduler/internal/domain"
	"time"
)

// NotificationItem struct that is saved into the DB
type NotificationItem struct {
	ID         string     `json:"id"`
	TelegramID string     `json:"telegram_id,omitempty"`
	Email      string     `json:"email,omitempty"`
	Message    string     `json:"message"`
	Via        domain.Via `json:"via"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	LastSent   *time.Time `json:"last_sent,omitempty"`
}

// CreateItemFromNotification creates a NotificationItem from a domain.Notification. It receives the transactionTi
func CreateItemFromNotification(notification domain.Notification) NotificationItem {
	notificationID := notification.ID
	if notificationID == "" {
		notificationID = uuid.NewString()
	}

	return NotificationItem{
		ID:         notificationID,
		TelegramID: notification.TelegramID,
		Email:      notification.Email,
		Message:    notification.Message,
		Via:        notification.Via,
		StartDate:  notification.StartDate,
		EndDate:    notification.EndDate,
	}
}

// ToNotification returns the NotificationItem as a domain.Notification
func (ni NotificationItem) ToNotification() domain.Notification {
	return domain.Notification{
		ID:         ni.ID,
		TelegramID: ni.TelegramID,
		Email:      ni.Email,
		Message:    ni.Message,
		Via:        ni.Via,
		StartDate:  ni.StartDate,
		EndDate:    ni.EndDate,
	}
}
