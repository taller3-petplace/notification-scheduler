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

func (fake *FakeDB) CreateNotifications(notifications []domain.Notification) error {
	for _, notification := range notifications {
		notificationItem := item.CreateItemFromNotification(notification)
		for _, hour := range notification.Hours {
			if !utils.ValidHour(hour) {
				return fmt.Errorf("error creating notifications: invalid key")
			}

			fake.db[hour] = append(fake.db[hour], notificationItem)
		}
	}

	return nil
}

func (fake *FakeDB) GetNotification(notificationID string) (*domain.Notification, error) {
	if fake.err != nil {
		return nil, fake.err
	}

	for _, notificationsPerHour := range fake.db {
		for _, notifItem := range notificationsPerHour {
			if notifItem.ID == notificationID {
				notification := notifItem.ToNotification()
				return &notification, nil
			}
		}
	}

	return nil, nil
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
