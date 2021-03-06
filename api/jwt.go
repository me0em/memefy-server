// Package api represents REST API and implements
// database abstraction, url routing and some jwt layer functions
package api

/*

TODO: refresh tokens system

TODO:
 - Who are the consumers of your tokens? Just internal?
 - Is it acceptable to have shared state?
	Can you put something opaque like a session id in your token and store
	the rest of the details inside your service?
 - What will you do if you need to quickly expire tokens that have been issued?
 -What is the scope of what can be accessed with one of your tokens?
	What could happen if one fell into the wrong hands?
*/

import (
	"errors"
	"memefy-server/config"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken allows to generate tokens (wow really?)
func GenerateToken(username string) string {
	// Create JWT token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := jwt.StandardClaims{}
	claims.Id = username
	// Expire in 3 days
	claims.ExpiresAt = time.Now().Add(time.Hour * 72).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		panic(err)
	}

	return tokenString
}

// Authorization checks a token for validation
// TODO: expired time
func Authorization(w http.ResponseWriter, r *http.Request) (string, float64, error) {
	// Parse token in request header
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {
		return "failed", 0, errors.New("expected valid bearer token")
	}

	token, _ := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return "failed", 0, errors.New("expected valid bearer token")
	}
	// IMHO not very smart move: wrap claims as interface :/
	userID := claims["jti"].(string)
	expTime := claims["exp"].(float64)

	return userID, expTime, nil
}

// TODO: дописать
// CheckRefreshToken checks refresh token for valid
func CheckRefreshToken(w http.ResponseWriter, r *http.Request) (string, error) {
	userID, _, err := Authorization(w, r)
	if err != nil {
		return "", err
	}
	if userID[len(userID)-8:] != "-refresh" {
		return "", errors.New("expected valid bearer token")
	}
	return userID, err
}
