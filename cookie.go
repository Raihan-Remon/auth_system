package main

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func TokenHandler(r *http.Request) (string, *Claims) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Fatal(err)
		}
		log.Fatal(err)

	}
	tokenString := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Fatal(err)

		}
		log.Fatal(err)
	}

	if !tkn.Valid {
		log.Fatal(err)
	}
	return tokenString, claims
}
