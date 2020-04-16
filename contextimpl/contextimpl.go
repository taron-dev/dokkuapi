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
const requestUserKey = key(3)
const requestAppKey = key(4)

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
		return nil, errors.New("Can't find github user in context")
	}
	return githubUser, nil
}

// DecorateWithUser decorate request with user
func DecorateWithUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), requestUserKey, user)
	return r.WithContext(ctx)
}

// GetUser returns user from request's context
func GetUser(ctx context.Context) (*model.User, error) {
	user, ok := ctx.Value(requestUserKey).(*model.User)
	if !ok {
		return nil, errors.New("Can't find user in context")
	}
	return user, nil
}

// DecorateWithApp decorate request with app
func DecorateWithApp(r *http.Request, app *model.Application) *http.Request {
	ctx := context.WithValue(r.Context(), requestAppKey, app)
	return r.WithContext(ctx)
}

// GetApp returns app from request's context
func GetApp(ctx context.Context) (*model.Application, error) {
	app, ok := ctx.Value(requestAppKey).(*model.Application)
	if !ok {
		return nil, errors.New("Can't find app in context")
	}
	return app, nil
}
