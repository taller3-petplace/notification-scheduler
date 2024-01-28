package service

import (
	"notification-scheduler/internal/domain"
)

type database interface {
}

type NotificationService struct {
	db database
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (ns *NotificationService) ScheduleNotification(notification domain.Notification) error {
	//TODO implement me
	panic("implement me")
}

func (ns *NotificationService) GetNotificationsByUserEmail(email string) ([]domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (ns *NotificationService) GetNotification(notificationID string) (domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (ns *NotificationService) DeleteNotification(notificationID string) error {
	//TODO implement me
	panic("implement me")
}
