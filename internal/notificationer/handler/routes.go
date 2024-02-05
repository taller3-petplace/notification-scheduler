package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"notification-scheduler/docs"
)

func (nh *NotificationHandler) RegisterRoutes(r *gin.Engine) {
	docs.SwaggerInfo.Title = "Swagger Notification Scheduler API"
	group := r.Group("/notifications", AppContextCreator())

	group.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
		return
	})
	group.POST("/notification", nh.ScheduleNotification)
	group.GET("/notification", nh.GetNotifications)
	group.GET("/notification/:notificationID", nh.GetNotificationData)
	group.PATCH("/notification/:notificationID", nh.UpdateNotification)
	group.DELETE("/notification/:notificationID", nh.DeleteNotification)
	group.POST("/email", nh.SendEmail)

	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
