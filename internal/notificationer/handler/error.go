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
	errInvalidUpdateRequest          = errors.New("error invalid update request body")
	errInvalidMail                   = errors.New("error invalid mail structure")
	errSendingEmail                  = errors.New("error sending email")
	errNotificationRequestValidation = errors.New("error notification request validation")
	errUpdateRequestValidation       = errors.New("error update notification request validation")
	errSchedulingNotification        = errors.New("error scheduling notification")
	errUserNotAllowed                = errors.New("error user not allowed")
	errFetchingUserNotifications     = errors.New("error fetching notifications")
	errFetchingNotification          = errors.New("error fetching notification")
	errUpdatingNotification          = errors.New("error updating")
	errMissingNotificationID         = errors.New("error missing notificationID")
	errDeletingNotification          = errors.New("error deleting notification")
)

var statusCodeByErr = map[error]int{
	errGettingAppContext:             http.StatusInternalServerError,
	errSchedulingNotification:        http.StatusInternalServerError,
	errFetchingUserNotifications:     http.StatusInternalServerError,
	errDeletingNotification:          http.StatusInternalServerError,
	errSendingEmail:                  http.StatusInternalServerError,
	errInvalidNotificationBody:       http.StatusBadRequest,
	errNotificationRequestValidation: http.StatusBadRequest,
	errMissingNotificationID:         http.StatusBadRequest,
	errInvalidMail:                   http.StatusBadRequest,
	errUpdateRequestValidation:       http.StatusBadRequest,
	errUserNotAllowed:                http.StatusUnauthorized,
}

func NewErrorResponse(err error) ErrorResponse {
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

	for errKey, statusCode := range statusCodeByErr {
		if errors.Is(err, errKey) {
			return ErrorResponse{
				StatusCode: statusCode,
				Message:    err.Error(),
			}
		}
	}

	return ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}
}
