package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("No Authorization header provided")
	}

	authStrings := strings.Split(authorization, " ")
	if len(authStrings) != 2 {
		return "", errors.New("Missing Authorization header fields")
	}

	if authStrings[0] != "Bearer" {
		return "", errors.New("Unknown/Unsuported Authorization method")
	}

	return authStrings[1], nil
}

func GetApiKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("No Authorization header provided")
	}

	authStrings := strings.Split(authorization, " ")
	if len(authStrings) != 2 {
		return "", errors.New("Missing Authorization header fields")
	}

	if authStrings[0] != "ApiKey" {
		return "", errors.New("Unknown/Unsuported Authorization method")
	}

	return authStrings[1], nil
}
