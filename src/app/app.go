package app

import "github.com/gin-gonic/gin"

type handler interface {
	RegisterRoutes(r *gin.Engine)
	ScheduleNotification(c *gin.Context)
	GetNotifications(c *gin.Context)
	GetNotificationData(c *gin.Context)
	UpdateNotification(c *gin.Context)
	DeleteNotification(c *gin.Context)
}

type App struct {
	NotificationHandler handler
}

// NewApp initializes all dependencies that App requires
func NewApp() *App {
	// DB

	// Service

	// Logger

	// Controller

	// App

	return &App{}
}
