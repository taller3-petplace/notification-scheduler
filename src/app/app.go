package app

import (
	"github.com/gin-gonic/gin"
	"notification-scheduler/internal/notificationer/db"
	"notification-scheduler/internal/notificationer/handler"
	"notification-scheduler/internal/notificationer/service"
)

type appHandler interface {
	RegisterRoutes(r *gin.Engine)
	ScheduleNotification(c *gin.Context)
	GetNotifications(c *gin.Context)
	GetNotificationData(c *gin.Context)
	UpdateNotification(c *gin.Context)
	DeleteNotification(c *gin.Context)
}

type App struct {
	NotificationHandler appHandler
}

// NewApp initializes all dependencies that App requires
func NewApp() *App {
	// DB
	appDB := db.NewFakeDB(nil)

	// Service
	notificationService := service.NewNotificationService(appDB)

	// Logger

	// Handler
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// App
	return &App{
		NotificationHandler: notificationHandler,
	}
}

func (a *App) RegisterRoutes(r *gin.Engine) {
	a.NotificationHandler.RegisterRoutes(r)
}
