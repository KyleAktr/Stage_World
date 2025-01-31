package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var ItemsPerPage = 6

func SetItemsPerPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	valueStr := r.URL.Query().Get("value")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		http.Error(w, "Valeur invalide", http.StatusBadRequest)
		return
	}

	ItemsPerPage = value
	fmt.Fprintf(w, "itemsPerPage mis à jour : %d", ItemsPerPage)
}

// LocalArtist représente la structure des artistes retournés par l'API
type LocalArtist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	CreationDate int      `json:"creationDate"`
	Members      []string `json:"members"`
}

// PageData contient les données pour le template
type PageData struct {
	Artists     []LocalArtist
	CurrentPage int
	TotalPages  int
}

// fetchArtistsFromAPI appelle votre API pour récupérer une liste d'artistes
func fetchArtistsFromAPI(apiURL string) ([]LocalArtist, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur lors de l'appel de l'API: %s", resp.Status)
	}

	var artists []LocalArtist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	// Afficher les données reçues pour debug
	fmt.Printf("Données reçues de l'API: %+v\n", artists)

	return artists, nil
}

// Handler gère la logique pour récupérer l'artiste, puis afficher les informations
func Handler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le numéro de page depuis les paramètres de l'URL
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Récupérer tous les artistes
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	allArtists, err := fetchArtistsFromAPI(apiURL)
	if err != nil {
		http.Error(w, "Impossible de récupérer les artistes depuis l'API", http.StatusInternalServerError)
		return
	}

	// Calculer la pagination
	totalPages := (len(allArtists) + ItemsPerPage - 1) / ItemsPerPage
	if page > totalPages {
		page = totalPages
	}

	// Sélectionner les artistes pour la page courante
	start := (page - 1) * ItemsPerPage
	end := start + ItemsPerPage
	if end > len(allArtists) {
		end = len(allArtists)
	}

	pageData := PageData{
		Artists:     allArtists[start:end],
		CurrentPage: page,
		TotalPages:  totalPages,
	}

	// Charger et exécuter le template
	tmpl, err := template.ParseFiles("./html/artistes.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}
}
