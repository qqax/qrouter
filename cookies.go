package qrouter

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

// GetCookie receives cookie from gin.Context or raise http.StatusForbidden error, returns cookie or error
func GetCookie(writer http.ResponseWriter, request *http.Request, name string) (string, error) {
	cookie, err := request.Cookie(name)
	if err != nil {
		log.Error().Err(err).Msg("getCookie error: " + err.Error())
		WriteJSONError(writer, err.Error(), http.StatusUnauthorized)
		return "", err
	}
	return cookie.Value, nil
}

// DeleteCookie receives cookie from gin.Context, deletes this cookie from gin.Context and returns its value or error.
func DeleteCookie(writer http.ResponseWriter, request *http.Request, name string) (string, error) {
	cookie, err := request.Cookie(name)
	if err != nil {
		return "", err
	} else {
		value := cookie.Value
		cookie.Value = ""
		cookie.MaxAge = -1
		http.SetCookie(writer, cookie)
		return value, nil
	}
}
