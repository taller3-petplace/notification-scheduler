package context

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestExtractJWTData(t *testing.T) {
	/*The JWT is created from here https://jwt.io/.
	The algorithm is HS256 and the secret 'ay harringui'. Contains the following data:
		{
		  "user_id": "69-abc",
		  "email": "larrycapija@testmail.com",
		  "telegram_id": "123"
		}
	*/

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjktYWJjIiwiZW1haWwiOiJsYXJyeWNhcGlq" +
		"YUB0ZXN0bWFpbC5jb20iLCJ0ZWxlZ3JhbV9pZCI6IjEyMyJ9.tddxCgzgHvCPHBHsakVod6fiN6C5Hf5t57OgpZHaKig"
	previousAlgorithm := os.Getenv(cryptographicAlgorithmEnvVar)
	previousSecret := os.Getenv(secretEnvVar)

	defer func() {
		_ = os.Setenv(cryptographicAlgorithmEnvVar, previousAlgorithm)
		_ = os.Setenv(secretEnvVar, previousSecret)
	}()

	err := os.Setenv(cryptographicAlgorithmEnvVar, "HS256")
	require.NoError(t, err)
	err = os.Setenv(secretEnvVar, "ay harringui")
	require.NoError(t, err)

	jwtUserData, err := extractDataFromJWT(tokenString)
	require.NoError(t, err)
	assert.Equal(t, "69-abc", jwtUserData.userID)
	assert.Equal(t, "larrycapija@testmail.com", jwtUserData.email)
	assert.Equal(t, "123", jwtUserData.telegramID)
}
