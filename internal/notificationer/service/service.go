package service

import (
	"notification-scheduler/internal/domain"
)

type searchFunction func(notification domain.Notification) bool

type database interface {
	CreateNotifications(notification domain.Notification) ([]domain.Notification, error)
	GetNotificationsByEmail(email string) ([]domain.Notification, error)
	GetNotification(notificationID string) (*domain.Notification, error)
	UpdateNotification(notification domain.Notification) error
	DeleteNotification(notificationID string) (bool, error)
	GetAll(currentHour string) []domain.Notification
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
	createdNotifications, err := ns.db.CreateNotifications(notification)
	if err != nil {
		return nil, newInternalError("ScheduleNotifications", err, "")
	}

	return createdNotifications, nil
}

// GetNotificationsByUserEmail searches all the notifications that have the given email
func (ns *NotificationService) GetNotificationsByUserEmail(email string) ([]domain.Notification, error) {
	operation := "GetNotificationsByUserEmail"
	notifications, err := ns.db.GetNotificationsByEmail(email)
	if err != nil {
		return nil, newInternalError(operation, err, "email: "+email)
	}

	return notifications, nil
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

// UpdateNotification updated the content of the given notification
func (ns *NotificationService) UpdateNotification(updatedNotification domain.Notification) error {
	operation := "UpdateNotification"
	err := ns.db.UpdateNotification(updatedNotification)
	if err != nil {
		return newInternalError(operation, err, "notificationID: "+updatedNotification.ID)
	}

	return nil
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

func (ns *NotificationService) GetAll(currentHour string) ([]domain.Notification, error) {
	return ns.db.GetAll(currentHour), nil
}
