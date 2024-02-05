package telegram

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"notification-scheduler/internal/domain"
	"notification-scheduler/internal/externalservices/telegram/internal/notification"
	"notification-scheduler/internal/internal/headers"
	"os"
	"time"
)

const (
	telegramSecretEnvVar     = "TELEGRAM_SECRET"
	telegramAccessCodeEnvVar = "TELEGRAM_ACCESS_CODE"
	url                      = "https://api.lnt.digital/telegram/notifications"
)

// Telegramer makes requests against Telegram Service
type Telegramer struct {
	clientHTTP http.Client
}

func NewTelegramer(client http.Client) *Telegramer {
	return &Telegramer{
		clientHTTP: client,
	}
}

// SendNotifications sends all the notifications to Telegram Service
// ToDo: send chunks of notifications
func (t *Telegramer) SendNotifications(notifications []domain.Notification) error {
	var telegramNotifications []notification.TelegramNotification
	for idx := range notifications {
		telegramNotifications = append(telegramNotifications, notification.NewTelegramNotification(notifications[idx]))
	}

	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		err = fmt.Errorf("%w: %v", errCreatingRequest, err)
		logrus.Errorf("%v", err)
		return err
	}

	accessToken, err := createAccessToken()
	if err != nil {
		logrus.Errorf("error creating telegram access token: %v", err)
		return fmt.Errorf("error creating token: %v", err)
	}

	request.Header.Add(headers.JWT, accessToken)
	response, err := t.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing SendNotifications: %v", err)
		return fmt.Errorf("%w: %v", errPerformingRequest, err)
	}

	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if response == nil {
		logrus.Error("nil response from telegram service")
		return errNilResponse
	}

	err = errPolicyFunc(response)
	if err != nil {
		logrus.Errorf("error from telegram service: %v", err)
		return err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("error from telegram service: %v", err)
		return fmt.Errorf("%w: %v", errUnmarshallingResponse, err)
	}
	logrus.Infof("Notification sent summary: %s", string(responseBody))

	return nil
}

// createAccessToken required token to make requests against Telegram Service
func createAccessToken() (string, error) {
	accessCode := os.Getenv(telegramAccessCodeEnvVar)
	if accessCode == "" {
		return "", fmt.Errorf("error access code is missing")
	}

	secret := os.Getenv(telegramSecretEnvVar)
	if secret == "" {
		return "", fmt.Errorf("error telegram secret is missing")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"access_code": accessCode,
			"exp":         time.Now().Add(2 * time.Minute).Unix(),
		})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
