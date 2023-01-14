package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const MAX_IMAGE_SIZE = 1024 * 1024

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	coll, _, ctx := DbConnection()
	tmpl, err := template.ParseFiles(path.Join("templates", "index.html"))
	if err != nil {
		panic(err)
	}
	if r.Method == "POST" {
		// w.Header().Set("Content-Type", "multipart/form-data;")
		if err := r.ParseMultipartForm(MAX_IMAGE_SIZE); err != nil {
			http.Error(w, "The uploaded file is too big. Please choose a file that's less than 1MB in size", http.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := getHash([]byte(r.FormValue("password")))
		gender := r.FormValue("gender")
		image, header, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}
		defer image.Close()

		ext := strings.Split(header.Filename, ".")[1]
		h := sha1.New()
		io.Copy(h, image)
		fileName := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
		fmt.Println(fileName)
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		path := filepath.Join(wd, "public", "images", fileName)
		nf, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer nf.Close()
		image.Seek(0, 0)
		io.Copy(nf, image)

		fmt.Printf("Name: %s\n Email: %s\n Password:%s\n Gender:%s\n Image:%s\n", name, email, password, gender, fileName)
		user := User{
			Name:     name,
			Email:    email,
			Password: password,
			Gender:   gender,
			Image:    fileName,
		}
		fmt.Println(user)
		insertResult, err := coll.InsertOne(ctx, &user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		// err = tmpl.Execute(w, nil)
		// if err != nil {
		// 	panic(err)
		// }
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	} else {
		err = tmpl.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	}

}
