package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/KyleAktr/Stage_World/api"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/index.html")
}

func ArtisteHandler(w http.ResponseWriter, r *http.Request) {
    // Log l'URL complète pour le débogage
    log.Printf("URL complète : %s\n", r.URL.Path)
    
    // Extraire l'ID de l'artiste de l'URL
    path := r.URL.Path
    parts := strings.Split(path, "/")
    
    // Log les parties de l'URL
    log.Printf("Parties de l'URL : %v\n", parts)
    
    // L'ID devrait être la dernière partie non-vide
    var id string
    for i := len(parts) - 1; i >= 0; i-- {
        if parts[i] != "" {
            id = parts[i]
            break
        }
    }
    
    log.Printf("ID extrait : %s\n", id)
    
    if id == "" {
        http.Error(w, "ID d'artiste manquant", http.StatusBadRequest)
        return
    }
    
    // Récupérer les détails de l'artiste
    artist, err := api.GetArtistByID(id)
    if err != nil {
        log.Printf("Erreur lors de la récupération de l'artiste %s : %v\n", id, err)
        http.Error(w, "Artiste non trouvé", http.StatusNotFound)
        return
    }
    
    // Préparer les données pour le template
    data := api.ArtistPageData{
        Artist: artist,
    }
    
    // Charger et exécuter le template
    tmpl := template.Must(template.ParseFiles("html/artiste.html"))
    if err := tmpl.Execute(w, data); err != nil {
        log.Printf("Erreur lors de l'exécution du template : %v\n", err)
        http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
        return
    }
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/contact.html")
}

func erreurHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/404.html")
}

func setupRoutes() {
    // Servir les fichiers statiques
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Routes des pages
    http.HandleFunc("/index", HomeHandler)
    http.HandleFunc("/404", erreurHandler)
    http.HandleFunc("/artiste/", ArtisteHandler)
    http.HandleFunc("/artistes", api.Handler)
    http.HandleFunc("/contact", ContactHandler)
}

func main() {
    setupRoutes()
    log.Println("Serveur démarré sur le port :8089")
    log.Fatal(http.ListenAndServe(":8089", nil))
}
