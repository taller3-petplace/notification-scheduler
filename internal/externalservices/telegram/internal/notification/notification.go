package notification

import "notification-scheduler/internal/domain"

type TelegramNotification struct {
	TelegramID string `json:"telegram_id"`
	Message    string `json:"message"`
}

func NewTelegramNotification(notif domain.Notification) TelegramNotification {
	return TelegramNotification{
		TelegramID: notif.TelegramID,
		Message:    notif.Message,
	}
}
