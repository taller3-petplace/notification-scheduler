package service

import (
	"errors"
	"fmt"
)

var (
	errNotificationNotFound      = errors.New("error notification not found")
	errNotificationAlreadyExists = errors.New("error notification already exists")
)

type serviceError struct {
	serviceOperation string
	err              error
	extraData        string
	alreadyExists    bool
	notFound         bool
	dbError          bool
}

func newInternalError(operation string, err error, extraData string) error {
	return serviceError{
		serviceOperation: operation,
		err:              err,
		extraData:        extraData,
		dbError:          true,
	}
}

func newNotificationNotFoundError(operation string, extraData string) error {
	return serviceError{
		serviceOperation: operation,
		err:              errNotificationNotFound,
		extraData:        extraData,
		notFound:         true,
	}
}

func newNotificationAlreadyExistsError(operation string, extraData string) error {
	return serviceError{
		serviceOperation: operation,
		err:              errNotificationAlreadyExists,
		extraData:        extraData,
		alreadyExists:    true,
	}
}

func (se serviceError) Error() string {
	if se.extraData != "" {
		return fmt.Sprintf("%v: %s - operation: %s", se.err, se.extraData, se.serviceOperation)
	}

	return fmt.Sprintf("%v - operation: %s", se.err, se.serviceOperation)
}

func (se serviceError) NotFound() bool {
	return se.notFound
}

func (se serviceError) InternalError() bool {
	return se.dbError
}

func (se serviceError) AlreadyExists() bool {
	return se.alreadyExists
}
