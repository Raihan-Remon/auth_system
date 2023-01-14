package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(path.Join("templates", "update.html"))
	if err != nil {
		panic(err)
	}
	var user User
	coll, _, ctx := DbConnection()

	if r.Method != "POST" {
		_, claims := TokenHandler(r)

		err = coll.FindOne(ctx, bson.M{"email": claims.Email}).Decode(&user)

		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(w, user)
		if err != nil {
			panic(err)
		}
	} else {
		_, claims := TokenHandler(r)

		if err := r.ParseForm(); err != nil {
			fmt.Printf("There was an error is parsing the form: %v", err)
			return
		}
		name := r.FormValue("name")
		password := r.FormValue("password")
		if password != "" {
			result, err := coll.UpdateOne(ctx, bson.M{"email": claims.Email}, bson.D{{"$set", bson.M{"name": name, "password": getHash([]byte(password))}}})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Updated: ", result.ModifiedCount)
		} else {
			result, err := coll.UpdateOne(ctx, bson.M{"email": claims.Email}, bson.D{{"$set", bson.M{"name": name}}})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Updated: ", result.ModifiedCount)
		}

		// err = tmpl.Execute(w, nil)
		// if err != nil {
		// 	panic(err)
		// }
		// TODO: IMPLEMENT REDIRECT
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	}
}
