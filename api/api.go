package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// LocalArtist représente la structure des artistes retournés par l'API
type LocalArtist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	CreationDate int      `json:"creationDate"`
	Members      []string `json:"members"`
}

// ArtistDetails représente les détails complets d'un artiste
type ArtistDetails struct {
	ID           int           `json:"id"`
	Image        string        `json:"image"`
	Name         string        `json:"name"`
	CreationDate int          `json:"creationDate"`
	Members      []string      `json:"members"`
	Locations    Locations     `json:"locations"`
	Concerts     ConcertDates  `json:"concertDates"`
	Relations    Relations     `json:"relations"`
}

// Locations représente les lieux de concert
type Locations struct {
	Locations []string `json:"locations"`
}

// ConcertDates représente les dates de concert
type ConcertDates struct {
	Dates []string `json:"dates"`
}

// Relations représente les relations entre les lieux et les dates
type Relations struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

// PageData contient les données pour le template
type PageData struct {
	Artists     []LocalArtist
	CurrentPage int
	TotalPages  int
}

// ArtistPageData contient les données pour la page d'un artiste
type ArtistPageData struct {
	Artist *ArtistDetails
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

	return artists, nil
}

// GetArtistByID récupère les détails d'un artiste spécifique
func GetArtistByID(id string) (*ArtistDetails, error) {
	// Afficher l'URL pour le débogage
	artistURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", id)
	log.Printf("Tentative de récupération de l'artiste à l'URL : %s\n", artistURL)
	
	resp, err := http.Get(artistURL)
	if err != nil {
		log.Printf("Erreur lors de la requête HTTP : %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status code non-OK reçu : %d\n", resp.StatusCode)
		return nil, fmt.Errorf("artiste non trouvé (status: %d)", resp.StatusCode)
	}

	var artist ArtistDetails
	if err := json.NewDecoder(resp.Body).Decode(&artist); err != nil {
		log.Printf("Erreur lors du décodage JSON : %v\n", err)
		return nil, err
	}

	// Récupérer les informations supplémentaires (relations)
	relationsURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%s", id)
	log.Printf("Tentative de récupération des relations à l'URL : %s\n", relationsURL)
	
	respRel, err := http.Get(relationsURL)
	if err != nil {
		log.Printf("Erreur lors de la requête des relations : %v\n", err)
		return nil, err
	}
	defer respRel.Body.Close()

	if respRel.StatusCode == http.StatusOK {
		if err := json.NewDecoder(respRel.Body).Decode(&artist.Relations); err != nil {
			log.Printf("Erreur lors du décodage des relations : %v\n", err)
			return nil, err
		}
	} else {
		log.Printf("Status code non-OK pour les relations : %d\n", respRel.StatusCode)
	}

	return &artist, nil
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
		log.Println("Erreur API Locale:", err)
		return
	}

	// Calculer la pagination
	const itemsPerPage = 6
	totalPages := (len(allArtists) + itemsPerPage - 1) / itemsPerPage
	if page > totalPages {
		page = totalPages
	}

	// Sélectionner les artistes pour la page courante
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
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
		log.Printf("Erreur template: %v\n", err)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		log.Printf("Erreur exécution template: %v\n", err)
		return
	}
}
