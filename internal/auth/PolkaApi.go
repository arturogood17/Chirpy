package auth

import (
	"net/http"
	"strings"

	"errors"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("No authorization body in headers")
	}
	t := strings.Split(auth, " ")
	if len(t) < 2 || t[0] != "ApiKey" {
		return "", errors.New("ApiKey malformed")
	}
	return t[1], nil
}
