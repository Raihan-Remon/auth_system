package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var jwtKey = []byte("secret_key")

// TODO: IMPLEMENT MODEL
type User struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Gender   string `json:"gender" bson:"gender"`
	Image    string `json:"image" bson:"image"`
}

type Claims struct {
	Email string `json:"email" bson:"email"`
	jwt.StandardClaims
}

func main() {
	_, client, ctx := DbConnection()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/login", LoginHandler)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	r.HandleFunc("/home", HomeHandler)
	r.HandleFunc("/update", UpdateHandler)
	r.HandleFunc("/delete", DeleteHandler)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("Server started at: http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal("Listen and serve error:", err)
	}

}
