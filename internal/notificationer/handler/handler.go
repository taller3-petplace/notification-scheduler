package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"notification-scheduler/internal/domain"
	"notification-scheduler/internal/externalservices/email"
	"notification-scheduler/internal/internal/context"
	"notification-scheduler/internal/notificationer/handler/internal/validator"
)

type servicer interface {
	ScheduleNotifications(notification domain.Notification) ([]domain.Notification, error)
	GetNotificationsByUserEmail(email string) ([]domain.Notification, error)
	GetNotification(notificationID string) (domain.Notification, error)
	UpdateNotification(notification domain.Notification) error
	DeleteNotification(notificationID string) error
}

type emailService interface {
	SendEmail(email email.Mail) error
}

type NotificationHandler struct {
	service     servicer
	emailClient emailService
}

func NewNotificationHandler(service servicer, emailClient emailService) *NotificationHandler {
	return &NotificationHandler{
		service:     service,
		emailClient: emailClient,
	}
}

// ScheduleNotification receives a domain.NotificationRequest, performs validations and if it's all OK then one notification per each
// specified hour is saved.
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
		response = append(response, domain.NewNotificationResponse(createdNotifications[idx]))
	}
	c.JSON(http.StatusCreated, response)
}

func (nh *NotificationHandler) GetNotifications(c *gin.Context) {
	appContext, err := context.GetAppContext(c.Request.Context())
	if err != nil {
		errResponse := NewErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	userIDParam := c.Param("userID")
	if userIDParam == "" || userIDParam != appContext.UserID {
		errResponse := NewErrorResponse(fmt.Errorf("%w: userID %s", errUserNotAllowed, appContext.UserID))
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

// GetNotificationData fetches notifications by user email
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

// UpdateNotification updates attributes of certain notification. The attributes that can be updated are: message and end date
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
	err = c.ShouldBindJSON(updateRequest)
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
//	@Summary		send mail
//	@Description	Send mail to given user
//	@Tags			Mail
//	@Accept			json
//	@Produce		json
//	@Param			mail	body		Mail	true	"mail info"
//	@Success		201		{object}	nil
//	@Failure		400,404	{object}	nil
//	@Router			/mail-service/send/ [post]
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
