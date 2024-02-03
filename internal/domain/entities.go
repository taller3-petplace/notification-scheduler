package domain

import "time"

// Notification structure that acts like a DTO. Its attributes are:
// + ID: identifier of the notification. Needed for the different types of operations. Is a UUID
//
// + TelegramID / Email: info needed to send a notification to one of these services
//
// + Message: message to be sent to the user
//
// + Via: can be Telegram, Mail or Both. The notification will be delivery to one of these services, or both
//
// + StartDate: when the notification is triggered
//
// + EndDate: when the notifications should stop. If none data was pass to this attribute, the notification never ends
//
// + Hours: hours of the day on which the notification should be sent
type Notification struct {
	ID         string     `json:"id"`
	TelegramID string     `json:"telegram_id,omitempty"`
	Email      string     `json:"email,omitempty"`
	Message    string     `json:"message"`
	Via        Via        `json:"via"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Hours      []string   `json:"hours"`
}
