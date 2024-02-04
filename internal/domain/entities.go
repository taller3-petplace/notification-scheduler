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

func Merge(notification Notification, update UpdateNotificationRequest) Notification {
	mergeResult := Notification{
		ID:         notification.ID,
		TelegramID: notification.TelegramID,
		Email:      notification.Email,
		Message:    notification.Message,
		Via:        notification.Via,
		StartDate:  notification.StartDate,
		EndDate:    notification.EndDate,
		Hours:      notification.Hours,
	}

	if notification.Message != update.Message {
		mergeResult.Message = update.Message
	}

	if notification.EndDate != nil && update.EndDate == nil {
		mergeResult.EndDate = nil
	} else if notification.EndDate != nil && update.EndDate != nil && !notification.EndDate.Equal(*update.EndDate) {
		mergeResult.EndDate = update.EndDate
	} else if notification.EndDate == nil && update.EndDate != nil {
		mergeResult.EndDate = update.EndDate
	}

	return mergeResult
}
