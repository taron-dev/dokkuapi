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
var jwtBlacklist []string

func hasValidToken(w http.ResponseWriter, r *http.Request) bool {
	if r.Header["Authorization"] != nil {
		reqToken := r.Header.Get("Authorization")

		if isBlacklisted(reqToken) {
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

// AddToBlacklist blacklist jwt
func AddToBlacklist(r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	jwtBlacklist = append(jwtBlacklist, reqToken)
}

func isBlacklisted(val string) bool {
	for _, item := range jwtBlacklist {
		if item == val {
			return true
		}
	}
	return false
}

// ExtractSub return sub from jwt
func ExtractSub(request *http.Request) string {
	reqToken := request.Header.Get("Authorization")
	reqToken = strings.Split(reqToken, "Bearer ")[1]
	token, _ := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Error parsing jwt")
		}
		return mySigningKey, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	return claims["sub"].(string)
}
