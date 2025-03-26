package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

var secret []byte = []byte(pass)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if len(pass) > 0 {
			var token string

			cookie, err := r.Cookie("token")

			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			token = cookie.Value

			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})

			if !jwtToken.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func generateToken() (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	token, err := jwtToken.SignedString(secret)

	if err != nil {
		return "", err
	}

	return token, err
}
