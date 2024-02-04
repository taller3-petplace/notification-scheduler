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
	ID         string
	TelegramID string
	Email      string
	Message    string
	Via        Via
	StartDate  time.Time
	EndDate    *time.Time
	Hours      []string
}
