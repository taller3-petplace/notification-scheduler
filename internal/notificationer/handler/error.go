package handler

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// serviceError interface for errors that come from the service
type serviceError interface {
	error
	NotFound() bool
	AlreadyExists() bool
	InternalError() bool
}

var (
	errGettingAppContext             = errors.New("error getting app")
	errInvalidNotificationBody       = errors.New("error invalid notification request body")
	errNotificationRequestValidation = errors.New("error notification request validation")
	errSchedulingNotification        = errors.New("error scheduling notification")
	errUserNotAllowed                = errors.New("error user not allowed")
	errFetchingUserNotifications     = errors.New("error fetching notifications")
	errFetchingNotification          = errors.New("error fetching notification")
	errMissingNotificationID         = errors.New("error missing notificationID")
	errDeletingNotification          = errors.New("error deleting notification")
)

var statusCodeByErr = map[error]int{
	errGettingAppContext:             http.StatusInternalServerError,
	errSchedulingNotification:        http.StatusInternalServerError,
	errFetchingUserNotifications:     http.StatusInternalServerError,
	errDeletingNotification:          http.StatusInternalServerError,
	errInvalidNotificationBody:       http.StatusBadRequest,
	errNotificationRequestValidation: http.StatusBadRequest,
	errMissingNotificationID:         http.StatusBadRequest,
	errUserNotAllowed:                http.StatusUnauthorized,
}

func NerErrorResponse(err error) ErrorResponse {
	var serviceErrorData serviceError
	isServiceError := errors.As(err, &serviceErrorData)
	if isServiceError && serviceErrorData.NotFound() {
		return ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    serviceErrorData.Error(),
		}
	}

	if isServiceError && serviceErrorData.InternalError() {
		return ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    serviceErrorData.Error(),
		}
	}

	errCode, ok := statusCodeByErr[err]
	if ok {
		return ErrorResponse{
			StatusCode: errCode,
			Message:    err.Error(),
		}
	}

	return ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}
}
