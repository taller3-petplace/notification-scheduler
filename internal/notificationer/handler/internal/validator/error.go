package validator

import "errors"

var (
	errInvalidMessage         = errors.New("error invalid message")
	errInvalidStartDate       = errors.New("error invalid start date")
	errInvalidEndDate         = errors.New("error invalid end date")
	errInvalidHour            = errors.New("error invalid hour")
	errInvalidVia             = errors.New("error invalid via")
	errMissingTelegramID      = errors.New("error missing telegramID")
	errMissingEmail           = errors.New("error missing telegramID")
	errMissingUserInformation = errors.New("error missing user information")
)
