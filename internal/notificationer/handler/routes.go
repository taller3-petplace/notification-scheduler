package handler

import "github.com/gin-gonic/gin"

func (nh *NotificationHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/notifications") // ToDo: add middleware. Licha

	group.POST("/notification", nh.ScheduleNotification) // telegram or front
	group.GET("/notification/:userID", nh.GetNotifications)
	group.GET("/notification/:notificationID", nh.GetNotificationData)
	group.PATCH("/notification/:notificationID", nh.UpdateNotification)
	group.DELETE("/notification/:notificationID", nh.DeleteNotification)
}