package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearerAuth := headers.Get("Authorization")
	if bearerAuth == "" {
		return "", fmt.Errorf("no authorization header exists")
	}
	return strings.TrimPrefix(bearerAuth, "Bearer "), nil
}
