package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

	"go.mongodb.org/mongo-driver/bson"
)

func GetHomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(path.Join("templates", "home.html"))
	if err != nil {
		panic(err)
	}
	var user User
	coll, _, ctx := DbConnection()

	_, claims := TokenHandler(r)
	err = coll.FindOne(ctx, bson.M{"email": claims.Email}).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
	err = tmpl.Execute(w, user)
	if err != nil {
		panic(err)
	}

}
