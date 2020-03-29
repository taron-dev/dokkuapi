package authentication

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/ondro2208/dokkuapi/logger"
	"net/http"
	"os"
	"strings"
	"time"
)

var mySigningKey = []byte(os.Getenv("JWT_TOKEN_SECRET"))

// IsAuthenticated verifies if request is authenticated
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

// GenerateJWT with userId included
func GenerateJWT(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
		"sub": userId,
	}

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.ErrorLogger.Fatalf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
