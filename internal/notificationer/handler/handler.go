package handler

import "github.com/gin-gonic/gin"

type servicer interface {
}

type NotificationHandler struct {
	service servicer
}

func NewNotificationHandler(service servicer) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

func (nh *NotificationHandler) ScheduleNotification(c *gin.Context) {
	panic("implement me!")
}

func (nh *NotificationHandler) GetNotifications(c *gin.Context) {
	panic("implement me!")
}

func (nh *NotificationHandler) GetNotificationData(c *gin.Context) {
	panic("implement me!")
}

func (nh *NotificationHandler) UpdateNotification(c *gin.Context) {
	panic("implement me!")
}

func (nh *NotificationHandler) DeleteNotification(c *gin.Context) {
	panic("implement me!")
}
