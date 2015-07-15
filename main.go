package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type IndexPage struct {
	Photos []string
}
type LoginPage struct {
	Body      string
	FirstName string
	LastName  string
	Email     string
}

func getPhotos() []string {
	photos := make([]string, 0)
	filepath.Walk("assets/images", func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		path = strings.Replace(path, "\\", "/", -1)
		photos = append(photos, path)
		return nil
	})
	return photos
}

func main() {
	http.Handle("/assets/",
		http.StripPrefix("/assets",
			http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		tpl, err := template.ParseFiles("assets/templates/index.gohtml")
		if err != nil {
			fmt.Println(err)
			http.Error(res, err.Error(), 500)
			return
		}
		err = tpl.Execute(res, IndexPage{
			Photos: getPhotos(),
		})
		if err != nil {
			fmt.Println(err)
			http.Error(res, err.Error(), 500)
			return
		}
	})
	http.HandleFunc("/admin", func(res http.ResponseWriter, req *http.Request) {

	})
	http.HandleFunc("/admin/login", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		tpl := template.New("assets.templates/login.gohtml")
		tpl, err := tpl.ParseFiles("assets/templates/login.gohtml")
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		err = tpl.Execute(res, LoginPage{
			FirstName: req.FormValue("firstName"),
			LastName:  req.FormValue("lastName"),
			Email:     req.FormValue("email"),
		})
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		http.SetCookie(res, &http.Cookie{
			Name: "loggedin", Value: "TRUE",
		})
	})
	http.HandleFunc("/admin/logout", func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name: "loggedin", Value: "",
		})
		http.Redirect(res, req, "/", 302)
	})

	http.ListenAndServe(":9000", nil)
}
