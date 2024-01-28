package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"notification-scheduler/internal/domain"
	"notification-scheduler/internal/internal/context"
	"notification-scheduler/internal/notificationer/handler/internal/validator"
)

type servicer interface {
	ScheduleNotification(notification domain.Notification) error
	GetNotificationsByUserEmail(email string) ([]domain.Notification, error)
	GetNotification(notificationID string) (domain.Notification, error)
	DeleteNotification(notificationID string) error
}

type NotificationHandler struct {
	service servicer
}

func NewNotificationHandler(service servicer) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

func (nh *NotificationHandler) ScheduleNotification(c *gin.Context) {
	appContext, err := context.GetAppContext(c)
	if err != nil {
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	var notificationRequest domain.NotificationRequest
	err = c.ShouldBindJSON(&notificationRequest)
	if err != nil {
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errInvalidNotificationBody, err))
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
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errNotificationRequestValidation, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notification := notificationRequest.ToNotification()
	err = nh.service.ScheduleNotification(notification)
	// ToDo: handle idempotent messages
	if err != nil {
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errSchedulingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notificationResponse := domain.NewNotificationResponse(notification)
	c.JSON(http.StatusCreated, notificationResponse)
}

func (nh *NotificationHandler) GetNotifications(c *gin.Context) {
	appContext, err := context.GetAppContext(c)
	if err != nil {
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	userIDParam := c.Param("userID")
	if userIDParam == "" || userIDParam != appContext.UserID {
		errResponse := NerErrorResponse(fmt.Errorf("%w: userID %s", errUserNotAllowed, appContext.UserID))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notifications, err := nh.service.GetNotificationsByUserEmail(appContext.Email)
	if err != nil {
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errFetchingUserNotifications, err))
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

func (nh *NotificationHandler) GetNotificationData(c *gin.Context) {
	appContext, err := context.GetAppContext(c)
	if err != nil {
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notificationID := c.Param("notificationID")
	if notificationID == "" {
		errResponse := NerErrorResponse(errMissingNotificationID)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notification, err := nh.service.GetNotification(notificationID)
	if err != nil {
		// ToDo: handle 404a
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errFetchingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if notification.Email != appContext.Email {
		errResponse := NerErrorResponse(fmt.Errorf("%w: userID %s", errUserNotAllowed, appContext.UserID))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	response := domain.NewNotificationResponse(notification)
	c.JSON(http.StatusOK, response)
}

func (nh *NotificationHandler) UpdateNotification(c *gin.Context) {
	panic("implement me!")
}

func (nh *NotificationHandler) DeleteNotification(c *gin.Context) {
	appContext, err := context.GetAppContext(c)
	if err != nil {
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errGettingAppContext, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	notificationID := c.Param("notificationID")
	if notificationID == "" {
		errResponse := NerErrorResponse(errMissingNotificationID)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	// Sanity check: only the user that creates the notification can delete it
	notification, err := nh.service.GetNotification(notificationID)
	if err != nil {
		// ToDo: handle 404
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errFetchingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	if notification.Email != appContext.Email {
		errResponse := NerErrorResponse(fmt.Errorf("%w: cannot delete notification, userID %s", errUserNotAllowed, appContext.UserID))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	err = nh.service.DeleteNotification(notificationID)
	if err != nil {
		errResponse := NerErrorResponse(fmt.Errorf("%w: %v", errDeletingNotification, err))
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, nil)
}
