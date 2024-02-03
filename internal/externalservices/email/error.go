package email

import "errors"

var (
	errCreatingSession = errors.New("error creating session")
	errSendingEmail    = errors.New("error sending email")
)
