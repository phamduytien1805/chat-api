package http_adapter

import (
	"errors"
	"net/http"
	"strings"
)

type contextKey string

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = contextKey("authorization_payload")
	authorizationRefreshKey = "refresh_token"
)

func setRfTokenCookie(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     authorizationRefreshKey,
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
	})

}

func getRfTokenFromCookie(r *http.Request) (string, error) {
	refreshCookie, err := r.Cookie(authorizationRefreshKey)
	if err != nil {
		return "", err
	}
	return refreshCookie.Value, nil
}

func parseAuthorizationHeader(header string) (authType, token string, err error) {
	fields := strings.Fields(header)
	if len(fields) != 2 {
		return "", "", errors.New("invalid authorization header format")
	}
	return fields[0], fields[1], nil
}
