package domain

import "time"

// Notification structure that is saved in the DB. Also, from it a response ir created when the user
// creates one for the first time. Its attributes are:
// + ID: identifier of the notification. Needed for the different types of operations
//
// + UserID: mail or telegram ID of the user
//
// + Message: message to be sent to the user
//
// + Via: can be Telegram or Mail, the notification will be delivery from one of them
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
	ID        int        `json:"id"`
	UserID    string     `json:"userID"`
	Message   string     `json:"message"`
	Via       via        `json:"via"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Frequency int        `json:"frequency"`
	LastSent  time.Time  `json:"last_sent"`
}
