package domain

import "time"

// Notification structure that is saved in the DB. Also, from it a response ir created when the user
// creates one for the first time. Its attributes are:
// + ID: identifier of the notification. Needed for the different types of operations
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
// + Frequency: its unit is in hours. ToDo: check this licha
//
// + LastSent: field that can be modified each time a notification is sent. Together with Frequency the service can know
// if it should send a new one or not
type Notification struct {
	ID         int        `json:"id"`
	TelegramID string     `json:"telegram_id,omitempty"`
	Email      string     `json:"email,omitempty"`
	Message    string     `json:"message"`
	Via        Via        `json:"via"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Frequency  int        `json:"frequency"`
	LastSent   time.Time  `json:"last_sent"`
}
