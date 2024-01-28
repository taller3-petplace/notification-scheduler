package domain

import "time"

// via service by which the notification is sent
type via string

const (
	Telegram via = "telegram"
	Mail     via = "mail"
)

type NotificationRequest struct {
	UserID    string     `json:"user_id" binding:"required"`
	Message   string     `json:"message" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required"`
	EndDate   *time.Time `json:"end_date"`
	Frequency int        `json:"frequency" binding:"required"`
	Times     int        `json:"times" binding:"required"`
	Via       via
}

func (nr NotificationRequest) ToNotification() Notification {
	return Notification{
		UserID:    nr.UserID,
		Message:   nr.Message,
		Via:       nr.Via,
		StartDate: nr.StartDate,
		EndDate:   nr.EndDate,
		Frequency: nr.Frequency,
	}
}

type NotificationResponse struct {
	ID        int        `json:"id"`
	Via       via        `json:"via"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

func NewNotificationResponse(notification Notification) NotificationResponse {
	return NotificationResponse{
		ID:        notification.ID,
		Via:       notification.Via,
		StartDate: notification.StartDate,
		EndDate:   notification.EndDate,
	}
}
