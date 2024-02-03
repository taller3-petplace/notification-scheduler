package domain

import (
	"notification-scheduler/internal/utils"
	"time"
)

// Via service by which the notification is sent
type Via string

const (
	Telegram Via = "telegram"
	Mail     Via = "mail"
	Both     Via = "both"
)

var validVias = []Via{
	Telegram,
	Mail,
	Both,
}

// ValidVia returns true if the given via is valid, otherwise false
func ValidVia(via Via) bool {
	return utils.Contains(validVias, via)
}

type NotificationRequest struct {
	TelegramID string     `json:"user_id"`
	Via        Via        `json:"via"`
	Message    string     `json:"message" binding:"required"`
	StartDate  time.Time  `json:"start_date" binding:"required"`
	EndDate    *time.Time `json:"end_date"`
	Hours      []string   `json:"hours" binding:"required"`
	Email      string
}

func (nr NotificationRequest) ToNotification() Notification {
	return Notification{
		TelegramID: nr.TelegramID,
		Email:      nr.Email,
		Message:    nr.Message,
		Via:        nr.Via,
		StartDate:  nr.StartDate,
		EndDate:    nr.EndDate,
		Hours:      nr.Hours,
	}
}

type NotificationResponse struct {
	ID        string     `json:"id"`
	Via       Via        `json:"via"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Hours     []string   `json:"hours"`
}

func NewNotificationResponse(notification Notification) NotificationResponse {
	return NotificationResponse{
		ID:        notification.ID,
		Via:       notification.Via,
		StartDate: notification.StartDate,
		EndDate:   notification.EndDate,
	}
}
