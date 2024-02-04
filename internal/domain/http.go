package domain

import (
	"encoding/json"
	"notification-scheduler/internal/utils"
	"strings"
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

// getViaFromString returns the matching via based on the given input. If the input does not match
// with any Via, is converted to one which might be invalid
func getViaFromString(input string) Via {
	input = strings.ToLower(input)
	switch input {
	case string(Telegram):
		return Telegram
	case string(Mail):
		return Mail
	case string(Both):
		return Both
	default:
		return Via(input)
	}
}

type NotificationRequest struct {
	TelegramID string     `json:"telegram_id"`
	Via        Via        `json:"via" binding:"required"`
	Message    string     `json:"message" binding:"required"`
	StartDate  time.Time  `json:"start_date" binding:"required"`
	EndDate    *time.Time `json:"end_date"`
	Hours      []string   `json:"hours" binding:"required"`
	Email      string
}

func (nr *NotificationRequest) UnmarshalJSON(rawData []byte) error {
	var requestData struct {
		TelegramID string     `json:"telegram_id"`
		Via        string     `json:"via"`
		Message    string     `json:"message"`
		StartDate  time.Time  `json:"start_date"`
		EndDate    *time.Time `json:"end_date"`
		Hours      []string   `json:"hours"`
	}

	err := json.Unmarshal(rawData, &requestData)
	if err != nil {
		return err
	}

	nr.TelegramID = requestData.TelegramID
	nr.Via = getViaFromString(requestData.Via)
	nr.Message = requestData.Message
	nr.StartDate = requestData.StartDate
	nr.EndDate = requestData.EndDate
	nr.Hours = requestData.Hours
	return nil
}

func (nr *NotificationRequest) ToNotification() Notification {
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
	Hour      string     `json:"hour"`
}

func NewNotificationResponse(notification Notification) NotificationResponse {
	return NotificationResponse{
		ID:        notification.ID,
		Via:       notification.Via,
		StartDate: notification.StartDate,
		EndDate:   notification.EndDate,
		Hour:      notification.Hours[0],
	}
}
