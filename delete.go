package main

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	coll, _, ctx := DbConnection()

	tokenString, claims := TokenHandler(r)
	result, err := coll.DeleteOne(ctx, bson.M{"email": claims.Email})
	if err != nil {
		panic(err)
	}
	fmt.Println(result.DeletedCount)
	http.SetCookie(w,
		&http.Cookie{
			Name:  "token",
			Value: "",
		})
	fmt.Println(tokenString)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}
