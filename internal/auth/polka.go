package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	s := headers.Get("Authorization")
	if s == "" {
		return "", errors.New("No authorization")
	}
	slicedString := strings.Split(s, " ")
	if len(slicedString) < 2 {
		return "", errors.New("malformed string")
	}
	if slicedString[0] != "ApiKey" {
		return "", errors.New("Not correct authorization")
	}
	apiKey := strings.TrimSpace(slicedString[1])

	return apiKey, nil
}
