package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (nh *NotificationHandler) RegisterRoutes(r *gin.Engine) {
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
}
