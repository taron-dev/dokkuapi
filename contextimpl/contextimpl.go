package contextimpl

import (
	"context"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ondro2208/dokkuapi/model"
	"net/http"
	"os"
	"strings"
)

type key int

const requestSubKey = key(1)
const requestGithubUserKey = key(2)

var mySigningKey = []byte(os.Getenv("JWT_TOKEN_SECRET"))

// DecorateWithSub decorate request with jwt's sub = user id
func DecorateWithSub(r *http.Request) *http.Request {
	reqToken := r.Header.Get("Authorization")
	reqToken = strings.Split(reqToken, "Bearer ")[1]
	token, _ := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Error parsing jwt")
		}
		return mySigningKey, nil
	})
	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	ctx := context.WithValue(r.Context(), requestSubKey, sub)
	return r.WithContext(ctx)

}

// GetSub returns user id from request's context
func GetSub(ctx context.Context) (string, error) {
	sub, ok := ctx.Value(requestSubKey).(string)
	if !ok {
		return "", errors.New("Can't find sub in context")
	}
	return sub, nil
}

// DecorateWithGithubUser decorate request with github user
func DecorateWithGithubUser(r *http.Request, githubUser *model.GithubUser) *http.Request {
	ctx := context.WithValue(r.Context(), requestGithubUserKey, githubUser)
	return r.WithContext(ctx)
}

// GetGithubUser returns github user from request's context
func GetGithubUser(ctx context.Context) (*model.GithubUser, error) {
	githubUser, ok := ctx.Value(requestGithubUserKey).(*model.GithubUser)
	if !ok {
		return nil, errors.New("Can't find github user id in context")
	}
	return githubUser, nil
}
