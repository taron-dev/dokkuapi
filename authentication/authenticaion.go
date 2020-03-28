package authentication

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var mySigningKey = []byte(os.Getenv("JWT_TOKEN_SECRET"))

func IsAuthenticated(endpointHandler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hasValidToken(w, r) {
			endpointHandler(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not Authorized"))
		}
	})
}

func hasValidToken(w http.ResponseWriter, r *http.Request) bool {
	if r.Header["Authorization"] != nil {
		reqToken := r.Header.Get("Authorization")
		reqToken = strings.Split(reqToken, "Bearer ")[1]
		token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Error parsing jwt")
			}
			return mySigningKey, nil
		})

		if err != nil {
			return false
		}

		if token.Valid {
			return true
		}

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}
