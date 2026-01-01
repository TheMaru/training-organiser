package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	noBearerErr := errors.New("No bearer")
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", noBearerErr
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return "", noBearerErr
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", noBearerErr
	}

	return parts[1], nil
}
