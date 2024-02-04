package db

import (
	"fmt"
	"notification-scheduler/internal/domain"
	"notification-scheduler/internal/notificationer/db/internal/item"
	"notification-scheduler/internal/utils"
)

type FakeDB struct {
	db  map[string][]item.NotificationItem
	err error
}

func NewFakeDB(err error) *FakeDB {
	db := make(map[string][]item.NotificationItem)
	return &FakeDB{
		db:  db,
		err: err,
	}
}

func (fake *FakeDB) CreateNotifications(notification domain.Notification) ([]domain.Notification, error) {
	var createdNotifications []domain.Notification
	for _, hour := range notification.Hours {
		if !utils.ValidHour(hour) {
			return nil, fmt.Errorf("error creating notifications: invalid key")
		}

		notificationItem := item.CreateItemFromNotification(notification)
		// Save notification
		fake.db[hour] = append(fake.db[hour], notificationItem)

		// Collect notifications. These notifications have a defined ID and hour
		createdNotification := notificationItem.ToNotification()
		createdNotification.Hours = []string{hour}
		createdNotifications = append(createdNotifications, createdNotification)
	}

	return createdNotifications, nil
}

func (fake *FakeDB) GetNotification(notificationID string) (*domain.Notification, error) {
	if fake.err != nil {
		return nil, fake.err
	}

	for hour, notificationsPerHour := range fake.db {
		for _, notifItem := range notificationsPerHour {
			if notifItem.ID == notificationID {
				notification := notifItem.ToNotification()
				notification.Hours = []string{hour}
				return &notification, nil
			}
		}
	}

	return nil, nil
}

func (fake *FakeDB) UpdateNotification(updatedNotification domain.Notification) error {
	if fake.err != nil {
		return fake.err
	}

	allNotifications, found := fake.db[updatedNotification.Hours[0]]
	if !found {
		return fmt.Errorf("key not found")
	}

	for idx := range allNotifications {
		if allNotifications[idx].ID == updatedNotification.ID {
			allNotifications[idx] = item.CreateItemFromNotification(updatedNotification)
			return nil
		}
	}

	return fmt.Errorf("error notification not found")
}

func (fake *FakeDB) DeleteNotification(notificationID string) (bool, error) {
	if fake.err != nil {
		return false, fake.err
	}

	deleted := false
	for key, notificationsPerHour := range fake.db {
		var notifItemsCopy []item.NotificationItem
		for idx := range notificationsPerHour {
			if notificationsPerHour[idx].ID == notificationID {
				deleted = true
				continue
			}
			notifItemsCopy = append(notifItemsCopy, notificationsPerHour[idx])
		}

		fake.db[key] = notifItemsCopy
	}

	return deleted, nil
}
