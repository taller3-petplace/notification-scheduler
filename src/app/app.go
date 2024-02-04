package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"notification-scheduler/internal/externalservices/email"
	"notification-scheduler/internal/notificationer/db"
	"notification-scheduler/internal/notificationer/handler"
	"notification-scheduler/internal/notificationer/service"
	"os"
)

type appHandler interface {
	RegisterRoutes(r *gin.Engine)
	ScheduleNotification(c *gin.Context)
	GetNotifications(c *gin.Context)
	GetNotificationData(c *gin.Context)
	UpdateNotification(c *gin.Context)
	DeleteNotification(c *gin.Context)
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

	// Handler
	notificationHandler := handler.NewNotificationHandler(notificationService, &session)

	// App
	return &App{
		NotificationHandler: notificationHandler,
	}, nil
}

func (a *App) RegisterRoutes(r *gin.Engine) {
	a.NotificationHandler.RegisterRoutes(r)
}
