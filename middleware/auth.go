package middleware

import "net/http"

// unfinished
func IsAuthorized(request *http.Request) bool {
	secretToken := "DICO"
	requestToken := request.Header.Get("Authorization")

	return secretToken == requestToken
}
