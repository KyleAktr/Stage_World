package main

import (
	"log"
	"net/http"

	"github.com/KyleAktr/Stage_World/api"
	"github.com/joho/godotenv"
)

func init() {
	// Charge les variables depuis le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Initialise l'API avec les variables d'environnement
	api.InitAPI()
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/index.html")
}
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/contact.html")
}
func ErreurHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/404.html")
}

func setupRoutes() {
	http.HandleFunc("/index", HomeHandler)
	http.HandleFunc("/404", ErreurHandler)
	http.HandleFunc("/artiste", api.ArtistHandler) // On appelle le handler artiste
	http.HandleFunc("/artistes", api.Handler)
	http.HandleFunc("/contact", ContactHandler)
}

func main() {
	// Gestion des fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Initialisation des routes
	setupRoutes()

	// Page 404 par défaut
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ErreurHandler(w, r)
	})

	http.ListenAndServe(":8089", nil)
}
