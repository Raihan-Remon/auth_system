package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	coll, _, ctx := DbConnection()
	// var user User
	var dbUser User
	tmpl, err := template.ParseFiles(path.Join("templates", "login.html"))
	if err != nil {
		panic(err)
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("There was an error is parsing the form: %v", err)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")

		err = coll.FindOne(ctx, bson.M{"email": email}).Decode(&dbUser)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dbUser)
		userPass := []byte(password)
		dbPass := []byte(dbUser.Password)

		passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

		if passErr != nil {
			log.Println(passErr)
			fmt.Println("Wrong Password!")
			return
		}

		fmt.Println("User Found!")

		expirationTime := time.Now().Add(time.Minute * 15)

		claims := &Claims{
			Email: dbUser.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			log.Fatal(err)
			return
		}

		http.SetCookie(w,
			&http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})
		fmt.Println(tokenString)
		// err = tmpl.Execute(w, nil)
		// if err != nil {
		// 	panic(err)
		// }
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	} else {
		err = tmpl.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	}
}
