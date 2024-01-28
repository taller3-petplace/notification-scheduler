package handler

import "github.com/gin-gonic/gin"

func (nh *NotificationHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/notifications", AppContextCreator())

	group.POST("/notification", nh.ScheduleNotification)
	group.GET("/notification/:userID", nh.GetNotifications)
	group.GET("/notification/:notificationID", nh.GetNotificationData)
	group.PATCH("/notification/:notificationID", nh.UpdateNotification)
	group.DELETE("/notification/:notificationID", nh.DeleteNotification)
}
