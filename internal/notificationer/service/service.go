package service

import (
	"notification-scheduler/internal/domain"
)

type searchFunction func(notification domain.Notification) bool

type database interface {
	CreateNotifications(notifications []domain.Notification) error
	GetNotification(notificationID string) (*domain.Notification, error)
	DeleteNotification(notificationID string) (bool, error)
}

type NotificationService struct {
	db database
}

func NewNotificationService(db database) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

// ScheduleNotifications creates the notifications. From one notification multiple can be created. This method
// contains all the logic to create the corresponding amount of notifications.
func (ns *NotificationService) ScheduleNotifications(notification domain.Notification) ([]domain.Notification, error) {

	return nil, nil
}

// GetNotificationsByUserEmail searches notifications based on the given search function. This function works as a filter
// in order to retrieve certain notifications
func (ns *NotificationService) GetNotificationsByUserEmail(email string) ([]domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

// GetNotification returns a single notification. If it does not exist, an error is returned
func (ns *NotificationService) GetNotification(notificationID string) (domain.Notification, error) {
	operation := "GetNotification"
	notification, err := ns.db.GetNotification(notificationID)
	if err != nil {
		return domain.Notification{}, newInternalError(operation, err, "notificationID: "+notificationID)
	}

	if notification == nil {
		return domain.Notification{}, newNotificationNotFoundError(operation, "notificationID: "+notificationID)
	}

	return *notification, err
}

// DeleteNotification deletes a single notification. If it does not exist, an error is returned
func (ns *NotificationService) DeleteNotification(notificationID string) error {
	operation := "DeleteNotification"
	deleted, err := ns.db.DeleteNotification(notificationID)
	if err != nil {
		return newInternalError(operation, err, "notificationID: "+notificationID)
	}

	if !deleted {
		return newNotificationNotFoundError(operation, "notificationID: "+notificationID)
	}

	return nil
}
