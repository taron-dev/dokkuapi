package authentication

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	log "github.com/ondro2208/dokkuapi/logger"
)

var mySigningKey = []byte(os.Getenv("JWT_TOKEN_SECRET"))
var jwtBlacklist []string

// HasValidToken validate if jwt in request is valid
func HasValidToken(w http.ResponseWriter, r *http.Request, blackList []string) bool {
	if r.Header["Authorization"] != nil {
		reqToken := r.Header.Get("Authorization")

		if isBlacklisted(reqToken, blackList) {
			return false
		}

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
func GenerateJWT(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
		"sub": userID,
	}

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.ErrorLogger.Fatalf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

// AddToBlacklist blacklist jwt
func AddToBlacklist(r *http.Request, blackList *[]string) {
	reqToken := r.Header.Get("Authorization")
	*blackList = append(*blackList, reqToken)
}

func isBlacklisted(val string, blackList []string) bool {
	for _, item := range blackList {
		if item == val {
			return true
		}
	}
	return false
}
