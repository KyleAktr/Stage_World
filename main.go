package main

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/index.html")
}
func ArtisteHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/artiste.html")
}
func ArtistesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/artistes.html")
}
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/contact.html")
}
func erreurHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/404.html")
}

func setupRoutes() {
	http.HandleFunc("/index", HomeHandler)
	http.HandleFunc("/404", erreurHandler)
	http.HandleFunc("/artiste", ArtisteHandler)
	http.HandleFunc("/artistes", ArtistesHandler)
	http.HandleFunc("/contact", ContactHandler)

}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	setupRoutes()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		erreurHandler(w, r)
	})

	http.ListenAndServe(":8089", nil)
}
