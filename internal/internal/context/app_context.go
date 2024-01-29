package context

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"notification-scheduler/internal/internal/headers"
	"os"
)

const (
	secretEnvVar                 = "secret"
	cryptographicAlgorithmEnvVar = "algorithm"
)

// AppContext context used by this app. It contains data, mainly from the user, that came in the request that can be use anywhere
type AppContext struct {
	TelegramRequest bool
	TelegramID      string
	UserID          string
	Email           string
}

// jwtData data that comes in the JWT
type jwtData struct {
	userID     string
	telegramID string
	email      string
}

type appContextKey struct{}

type appContextValue struct {
	Context AppContext
}

func NewAppContext(request *http.Request) (context.Context, error) {
	requestFromTelegram := request.Header.Get(headers.Telegram) == "true"
	appContext := AppContext{
		TelegramRequest: requestFromTelegram,
	}

	if !requestFromTelegram {
		tokenString := request.Header.Get(headers.JWT)
		tokenData, err := extractDataFromJWT(tokenString)
		if err != nil {
			return nil, fmt.Errorf("error extracting data from JWT: %v", err)
		}
		appContext.UserID = tokenData.userID
		appContext.Email = tokenData.email
		appContext.TelegramID = tokenData.telegramID
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

// extractDataFromJWT extracts all the data that the JWT contains. In order to do that it uses some
// environment variables. An error is returned if for some reason, the data cannot be extracted
func extractDataFromJWT(tokenString string) (*jwtData, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("error token is missing")
	}

	secretKey := os.Getenv(secretEnvVar)
	if secretKey == "" {
		return nil, fmt.Errorf("error secret key is missing")
	}

	algorithm := os.Getenv(cryptographicAlgorithmEnvVar)
	if algorithm == "" {
		return nil, fmt.Errorf("error algorithm is missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Sanity check: the algorithm must be the same as the one in the env var
		if token.Method.Alg() != algorithm {
			return nil, fmt.Errorf("error unexpected signing method: %s", token.Method.Alg())
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing JWT: %v", err)
	}

	// Access the claims if the signature is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &jwtData{
			userID:     claims["user_id"].(string),
			email:      claims["email"].(string),
			telegramID: claims["telegram_id"].(string),
		}, nil
	}

	return nil, fmt.Errorf("error invalid JWT")
}
