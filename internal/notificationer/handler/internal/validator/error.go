package validator

import "errors"

var (
	errInvalidMessage         = errors.New("error invalid message")
	errInvalidStartDate       = errors.New("error invalid start date")
	errInvalidEndDate         = errors.New("error invalid end date")
	errInvalidFrequency       = errors.New("error invalid frequency")
	errInvalidTimesPerDay     = errors.New("error invalid timer per day")
	errInvalidCombination     = errors.New("error invalid combination of parameters")
	errInvalidVia             = errors.New("error invalid via")
	errMissingTelegramID      = errors.New("error missing telegramID")
	errMissingEmail           = errors.New("error missing telegramID")
	errMissingUserInformation = errors.New("error missing user information")
)
