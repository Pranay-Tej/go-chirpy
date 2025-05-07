package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetPolkaApiKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("authorization header not found")
	}

	split := strings.Split(authorization, " ")
	if len(split) != 2 || strings.ToLower(split[0]) != "apikey" {
		return "", errors.New("invalid api key")
	}
	return split[1], nil
}
