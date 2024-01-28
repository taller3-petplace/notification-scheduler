package validator

import (
	"fmt"
	"notification-scheduler/internal/domain"
	"time"
)

// ValidateNotificationRequest validates the given notification request. The following checks are performed:
// + Message must be at least of length 5
// + StartDate and EndDate must be from now on, not from the past
// + At least one of them, Frequency or TimesPerDay, must be greater than zero
// + Via must be a valid one. Actually only Telegram, Mail or Both are valid
// + If via is 'telegram', the notification must contain the telegramID of the user
// + If via is 'mail', the notification must contain the email of the user
// + If via is 'both', the notification must contain the email and telegramId of the user
func ValidateNotificationRequest(notification domain.NotificationRequest) error {
	currentTime := time.Now()

	if len(notification.Message) < 5 {
		return fmt.Errorf("%w: must be of length at least 5", errInvalidMessage)
	}

	if notification.StartDate.Before(currentTime) {
		return fmt.Errorf("%w: date from the past", errInvalidStartDate)
	}

	if notification.EndDate != nil && notification.EndDate.Before(currentTime) {
		return fmt.Errorf("%w: date from the past", errInvalidEndDate)
	}

	if notification.Frequency < 0 {
		return fmt.Errorf("%w: negative frequency %v. Must be greater or equal to zero", errInvalidFrequency, notification.Frequency)
	}

	if notification.TimesPerDay < 0 {
		return fmt.Errorf("%w: negative times per day %v. Must be greater or equal to zero", errInvalidTimesPerDay, notification.TimesPerDay)
	}

	if notification.TimesPerDay == 0 && notification.Frequency == 0 {
		return fmt.Errorf("%w: times per day and frequency are zero. At least one of them has to be greater than zero", errInvalidCombination)
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
