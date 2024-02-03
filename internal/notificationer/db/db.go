package db

import (
	"github.com/guregu/dynamo"
	"notification-scheduler/internal/domain"
)

type Persistor struct {
	db *dynamo.DB
}

func NewPersistor() *Persistor {
	return &Persistor{}
}

// BatchInsert inserts multiple notifications
func (p *Persistor) BatchInsert(notifications []domain.Notification) error {
	return nil
}

// BatchGet inserts multiple notifications
func (p *Persistor) BatchGet(notifications []domain.Notification) ([]domain.Notification, error) {
	return nil, nil
}
