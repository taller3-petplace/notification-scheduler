package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"notification-scheduler/internal/domain"
	"notification-scheduler/internal/externalservices/email"
	"notification-scheduler/internal/internal/context"
	"notification-scheduler/internal/notificationer/handler/internal/validator"
	"time"
)

type servicer interface {
	ScheduleNotifications(notification domain.Notification) ([]domain.Notification, error)
	GetNotificationsByUserEmail(email string) ([]domain.Notification, error)
	GetNotification(notificationID string) (domain.Notification, error)
	UpdateNotification(notification domain.Notification) error
	DeleteNotification(notificationID string) error
	GetAll(hour string) ([]domain.Notification, error)
}

type emailService interface {
	SendEmail(email email.Mail) error
}

type telegramService interface {
	SendNotifications(notifications []domain.Notification) error
}

type NotificationHandler struct {
	service     servicer
	emailClient emailService
	telegramer  telegramService
}

func NewNotificationHandler(service servicer, emailClient emailService, telegramer telegramService) *NotificationHandler {
	return &NotificationHandler{
		service:     service,
		emailClient: emailClient,
		telegramer:  telegramer,
	}
}

// ScheduleNotification godoc
//
//	@Summary		Schedules notifications
//	@Description	Receives a domain.NotificationRequest, performs validations and if it's all OK then one notification per each specified hour is saved.
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string						true	"jwt"
//	@Param			X-Telegram-App		header		string						false	"true if the request comes from telegram service, otherwise false"
//	@Param			NotificationRequest	body		domain.NotificationRequest	true	"info about the notification to create"
//	@Success		201					{object}	[]domain.NotificationResponse
//	@Failure		400,404				{object}	nil
//	@Router			/notifications/notification [post]
func (nh *NotificationHandler) ScheduleNotification(c *gin.Context) {
	appContext, err := context.GetAppContext(c.Request.Context())
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	var notificationRequest domain.NotificationRequest
	err = c.ShouldBindJSON(&notificationRequest)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errInvalidNotificationBody, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	// If the request came from the frontend we have to set the IDs here
	if !appContext.TelegramRequest {
		notificationRequest.TelegramID = appContext.TelegramID
		notificationRequest.Email = appContext.Email
	} else {
		notificationRequest.Via = domain.Telegram
	}

	err = validator.ValidateNotificationRequest(notificationRequest)
	if err != nil {
		a := fmt.Errorf("%w: %v", errNotificationRequestValidation, err)
		errResponse := NewErrorResponse(a)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notification := notificationRequest.ToNotification()
	createdNotifications, err := nh.service.ScheduleNotifications(notification)
	var serviceErrorContext serviceError
	if errors.As(err, &serviceErrorContext) && serviceErrorContext.AlreadyExists() {
		c.JSON(http.StatusOK, domain.NewNotificationResponse(notification))
		return
	}

	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errSchedulingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	var response []domain.NotificationResponse
	for idx := range createdNotifications {
		notificationResponse := domain.NewNotificationResponse(createdNotifications[idx])
		notificationResponse.HideMessage()
		response = append(response, notificationResponse)
	}
	c.JSON(http.StatusCreated, response)
}

// GetNotifications godoc
//
//	@Summary		Search all notifications by user email
//	@Description	Returns all the notifications of the given user
//
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"jwt data, must contain the email of the user"
//	@Success		200				{object}	[]domain.NotificationResponse
//	@Failure		400,401,404		{object}	ErrorResponse
//	@Router			/notifications/notification [get]
func (nh *NotificationHandler) GetNotifications(c *gin.Context) {
	appContext, err := context.GetAppContext(c.Request.Context())
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notifications, err := nh.service.GetNotificationsByUserEmail(appContext.Email)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errFetchingUserNotifications, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if len(notifications) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	var response []domain.NotificationResponse
	for idx := range notifications {
		response = append(response, domain.NewNotificationResponse(notifications[idx]))
	}

	c.JSON(http.StatusOK, response)
}

// GetNotificationData godoc
//
//	@Summary		Fetches notification by ID
//	@Description	Fetches notification by ID
//
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"jwt data"
//	@Param			notificationID	path		string	true	"id of the notification"
//	@Success		200				{object}	domain.NotificationResponse
//	@Failure		400,401,404		{object}	ErrorResponse
//	@Router			/notifications/notification/{notificationID} [get]
func (nh *NotificationHandler) GetNotificationData(c *gin.Context) {
	appContext, err := context.GetAppContext(c.Request.Context())
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notificationID := c.Param("notificationID")
	if notificationID == "" {
		errResponse := NewErrorResponse(errMissingNotificationID)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notification, err := nh.service.GetNotification(notificationID)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %w", errFetchingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if notification.Email != appContext.Email {
		errResponse := NewErrorResponse(fmt.Errorf("%w: userID %s", errUserNotAllowed, appContext.UserID))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	response := domain.NewNotificationResponse(notification)
	c.JSON(http.StatusOK, response)
}

// UpdateNotification godoc
//
//	@Summary		Updates a notification
//	@Description	Updates attributes of certain notification. The attributes that can be updated are: message and end date
//
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"jwt data"
//	@Param			notificationID	path		string								true	"id of the notification"
//	@Param			UpdateRequest	body		domain.UpdateNotificationRequest	true	"Fields to update"
//	@Success		200				{object}	nil
//	@Failure		400,401,404		{object}	ErrorResponse
//	@Router			/notifications/notification/{notificationID} [patch]
func (nh *NotificationHandler) UpdateNotification(c *gin.Context) {
	appContext, err := context.GetAppContext(c.Request.Context())
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if appContext.TelegramRequest {
		errResponse := NewErrorResponse(fmt.Errorf("requests from Telegram are not allowed"))
		c.JSON(http.StatusForbidden, errResponse)
		return
	}

	var updateRequest domain.UpdateNotificationRequest
	err = c.ShouldBindJSON(&updateRequest)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errInvalidUpdateRequest, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	err = validator.ValidateUpdateRequest(updateRequest)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errUpdateRequestValidation, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notificationID := c.Param("notificationID")
	if notificationID == "" {
		errResponse := NewErrorResponse(errMissingNotificationID)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notification, err := nh.service.GetNotification(notificationID)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %w", errFetchingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	// Sanity check: the notification must belong to the user
	if notification.Email != appContext.Email {
		errResponse := NewErrorResponse(fmt.Errorf("%w: cannot update notification, userID %s", errUserNotAllowed, appContext.UserID))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	updatedNotification := domain.Merge(notification, updateRequest)
	err = nh.service.UpdateNotification(updatedNotification)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %w", errUpdatingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// DeleteNotification godoc
//
//	@Summary		Deletes a notification
//	@Description	If exists, deletes the notification with the given notificationID. This action is triggered by the users, if a notification reaches the end date nothing happens
//
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"jwt data"
//	@Param			notificationID	path		string	true	"id of the notification"
//	@Success		200				{object}	nil
//	@Failure		400,401,404		{object}	ErrorResponse
//	@Router			/notifications/notification/{notificationID} [delete]
func (nh *NotificationHandler) DeleteNotification(c *gin.Context) {
	appContext, err := context.GetAppContext(c.Request.Context())
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notificationID := c.Param("notificationID")
	if notificationID == "" {
		errResponse := NewErrorResponse(errMissingNotificationID)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	// Sanity check: only the user that creates the notification can delete it
	notification, err := nh.service.GetNotification(notificationID)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %w", errFetchingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if notification.Email != appContext.Email {
		errResponse := NewErrorResponse(fmt.Errorf("%w: cannot delete notification, userID %s", errUserNotAllowed, appContext.UserID))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	err = nh.service.DeleteNotification(notificationID)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %w", errDeletingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// SendEmail godoc
//
//	@Summary		Send mail
//	@Description	Send mail to given user
//	@Tags			Mail
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string		true	"jwt data"
//	@Param			mail			body		email.Mail	true	"mail info"
//	@Success		201				{object}	nil
//	@Failure		400,404			{object}	ErrorResponse
//	@Router			/notifications/email [post]
func (nh *NotificationHandler) SendEmail(c *gin.Context) {
	var mail email.Mail
	err := c.ShouldBindJSON(&mail)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %w", errInvalidMail, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	err = nh.emailClient.SendEmail(mail)
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %w", errSendingEmail, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// TriggerNotifications godoc
//
//	@Summary		sends notifications
//	@Description	Sends notifications to all users that have scheduled one for the hour of this request
//	@Tags			Notification
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"jwt data"
//	@Success		200				{object}	nil
//	@Success		204				{object}	nil
//	@Failure		400,404			{object}	ErrorResponse
//	@Router			/notifications/trigger [post]
func (nh *NotificationHandler) TriggerNotifications(c *gin.Context) {
	currentHour := time.Now().Hour()
	notifications, err := nh.service.GetAll(fmt.Sprint(currentHour))
	if err != nil {
		a := fmt.Errorf("error searching all notifications for given hour %d: %v", currentHour, err)
		errResponse := NewErrorResponse(a)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if len(notifications) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	for idx := range notifications {
		notification := notifications[idx]
		if notification.Via == domain.Mail || notification.Via == domain.Both {
			mail := email.Mail{
				To:      notification.Email,
				Subject: "Scheduled notification",
				Body:    notification.Message,
			}
			err = nh.emailClient.SendEmail(mail)
			if err != nil {
				logrus.Errorf("error sending mail: %v", err)
			}
			continue
		}

		if notification.Via == domain.Telegram || notification.Via == domain.Both {
			// ToDo: refactor. Licha
			err = nh.telegramer.SendNotifications([]domain.Notification{notification})
			if err != nil {
				logrus.Errorf("error sending telegram: %v", err)
			}
			continue
		}
	}

	c.JSON(http.StatusOK, nil)
}
