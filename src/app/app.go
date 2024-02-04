package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"notification-scheduler/internal/domain"
	"notification-scheduler/internal/externalservices/email"
	"notification-scheduler/internal/externalservices/telegram"
	"notification-scheduler/internal/notificationer/db"
	"notification-scheduler/internal/notificationer/handler"
	"notification-scheduler/internal/notificationer/service"
	"os"
	"time"
)

type appHandler interface {
	RegisterRoutes(r *gin.Engine)
	ScheduleNotification(c *gin.Context)
	GetNotifications(c *gin.Context)
	GetNotificationData(c *gin.Context)
	UpdateNotification(c *gin.Context)
	DeleteNotification(c *gin.Context)
}

type telegramHandler interface {
	SendNotifications(notifications []domain.Notification) error
}

func loadEmailConfig() (*email.EmailConfig, error) {

	region := os.Getenv("MAIL_REGION")
	if region == "" {
		return nil, errors.New("missing region")
	}
	secretKey := os.Getenv("MAIL_SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("missing secret key")
	}
	accessKey := os.Getenv("MAIL_ACCESS_KEY")
	if accessKey == "" {
		return nil, errors.New("missing access key")
	}
	from := os.Getenv("MAIL_FROM")
	if from == "" {
		return nil, errors.New("missing from")
	}

	return &email.EmailConfig{
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		From:      from,
	}, nil
}

type App struct {
	NotificationHandler appHandler
	Telegramer          telegramHandler
}

// NewApp initializes all dependencies that App requires
func NewApp() (*App, error) {
	// DB
	appDB := db.NewFakeDB(nil)

	// Service
	notificationService := service.NewNotificationService(appDB)

	// Aws Client
	emailConfig, err := loadEmailConfig()
	if err != nil {
		return nil, err
	}
	session := email.NewAwsSession(emailConfig)
	err = session.Connect()
	if err != nil {
		return nil, err
	}

	// Telegramer
	client := http.Client{Timeout: 5 * time.Second}
	telegramer := telegram.NewTelegramer(client)

	// Handler
	notificationHandler := handler.NewNotificationHandler(notificationService, &session)

	// App
	return &App{
		NotificationHandler: notificationHandler,
		Telegramer:          telegramer,
	}, nil
}

func (a *App) RegisterRoutes(r *gin.Engine) {
	a.NotificationHandler.RegisterRoutes(r)
}
