package validator

import (
	"fmt"
	"notification-scheduler/internal/domain"
	"notification-scheduler/internal/utils"
	"time"
)

// ValidateNotificationRequest validates the given notification request. The following checks are performed:
// + Message must be at least of length 5
// + StartDate and EndDate must be from now on, not from the past
// + The hours must be on the hour or thirty. Their range go from 0 to 23
// + Via must be a valid one. Actually only Telegram, Mail or Both are valid
// + If via is 'telegram', the notification must contain the telegramID of the user
// + If via is 'mail', the notification must contain the email of the user
// + If via is 'both', the notification must contain the email and telegramId of the user
func ValidateNotificationRequest(notification domain.NotificationRequest) error {
	//currentTime := time.Now()

	if len(notification.Message) < 5 {
		return fmt.Errorf("%w: must be of length at least 5", errInvalidMessage)
	}

	//if notification.StartDate.Before(currentTime) {
	//	return fmt.Errorf("%w: date from the past", errInvalidStartDate)
	//}
	//
	//if notification.EndDate != nil && notification.EndDate.Before(currentTime) {
	//	return fmt.Errorf("%w: date from the past", errInvalidEndDate)
	//}

	hoursSet := make(map[string]bool)
	for _, hour := range notification.Hours {
		if !utils.ValidHour(hour) {
			return fmt.Errorf("%w: hours must be o'clock or 30, and range from 0 to 23. Given: %s", errInvalidHour, hour)
		}
		if hoursSet[hour] {
			return fmt.Errorf("%w: %s", errRepeatedHour, hour)
		}
		hoursSet[hour] = true
	}

	if !domain.ValidVia(notification.Via) {
		return fmt.Errorf("%w: %s", errInvalidVia, notification.Via)
	}

	if notification.Via == domain.Telegram && notification.TelegramID == "" {
		return errMissingTelegramID
	}

	if notification.Via == domain.Mail && notification.Email == "" {
		return errMissingEmail
	}

	if notification.Via == domain.Both && (notification.Email == "" || notification.TelegramID == "") {
		return fmt.Errorf(
			"%w: email or telegramID is missing. Got email: %s - telegramID: %s",
			errMissingUserInformation,
			notification.Email,
			notification.TelegramID,
		)
	}

	return nil
}

// ValidateUpdateRequest validates the given update notification request. The following checks are performed:
// + Message must be at least of length 5
// + EndDate must be from now on, not from the past
// + At least one attribute must be updated
func ValidateUpdateRequest(notification domain.UpdateNotificationRequest) error {
	if notification.Message == "" && notification.EndDate == nil {
		return errNothingToUpdate
	}

	if len(notification.Message) < 5 {
		return fmt.Errorf("%w: must be of length at least 5", errInvalidMessage)
	}

	if notification.EndDate != nil && notification.EndDate.Before(time.Now()) {
		return fmt.Errorf("%w: date from the past", errInvalidEndDate)
	}

	return nil
}
