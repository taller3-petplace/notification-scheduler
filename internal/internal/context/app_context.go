package context

import (
	"context"
	"fmt"
	"net/http"
	"notification-scheduler/internal/internal/headers"
)

type AppContext struct {
	TelegramRequest bool
	UserID          string
	Email           string
}

type appContextKey struct{}

type appContextValue struct {
	Context AppContext
}

func NewAppContext(request *http.Request) (context.Context, error) {
	requestFromTelegram := request.Header.Get(headers.Telegram) == "true"
	token := request.Header.Get(headers.JWT)

	if token == "" {
		return nil, fmt.Errorf("missing jwt")
	}

	appContext := AppContext{
		TelegramRequest: requestFromTelegram,
		UserID:          "2", // ToDo: refactor this. Licha
	}

	return context.WithValue(
		request.Context(),
		appContextKey{},
		appContextValue{
			appContext,
		},
	), nil
}

// GetAppContext from the given context extracts the AppContext that should have been added by the middleware
func GetAppContext(ctx context.Context) (AppContext, error) {
	if ctx == nil {
		return AppContext{}, errNilContext
	}

	contextValue := ctx.Value(appContextKey{})
	if contextValue == nil {
		return AppContext{}, errMissingAppContext
	}

	appContext := contextValue.(appContextValue)
	return appContext.Context, nil
}
