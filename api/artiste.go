package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type ArtistPageData struct {
	Name            string
	SpotifyArtistID string
	CreationDate    int
	Members         []string
	Image           string
	Concerts        []Concert
}

type Concert struct {
	Location string
	Date     string
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// fetchConcerts récupère les dates de concert pour un artiste donné
func fetchConcerts(artistID int) (Relations, error) {
	var relations Relations
	relationsURL := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artistID)

	resp, err := http.Get(relationsURL)
	if err != nil {
		return relations, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return relations, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&relations)
	if err != nil {
		return relations, fmt.Errorf("error decoding response: %v", err)
	}

	return relations, nil
}

// Handler pour la page artiste
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("name")
	if artistName == "" {
		http.Error(w, "Artist name is required", http.StatusBadRequest)
		return
	}

	// Récupérer les informations de l'artiste depuis l'API
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	artists, err := fetchArtistsFromAPI(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch artists data", http.StatusInternalServerError)
		return
	}

	// Trouver l'artiste correspondant
	var artistData LocalArtist
	found := false
	for _, artist := range artists {
		if artist.Name == artistName {
			artistData = artist
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	// Récupérer les concerts
	relations, err := fetchConcerts(artistData.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch concert data: %v", err), http.StatusInternalServerError)
		return
	}

	// Créer la liste des concerts
	var concerts []Concert
	for location, dates := range relations.DatesLocations {
		for _, date := range dates {
			concerts = append(concerts, Concert{
				Location: location,
				Date:     date,
			})
		}
	}

	// Récupération du token Spotify
	token, err := GetSpotifyToken()
	if err != nil {
		http.Error(w, "Failed to get Spotify token", http.StatusInternalServerError)
		return
	}

	// Recherchez l'ID Spotify de l'artiste
	spotifyID, err := SearchArtist(artistName, token)
	if err != nil {
		http.Error(w, "Artist not found on Spotify", http.StatusNotFound)
		return
	}

	// Préparer les données pour le template
	data := ArtistPageData{
		Name:            artistData.Name,
		SpotifyArtistID: spotifyID,
		CreationDate:    artistData.CreationDate,
		Members:         artistData.Members,
		Image:           artistData.Image,
		Concerts:        concerts,
	}

	tmpl, err := template.ParseFiles("./html/artiste.html")
	if err != nil {
		http.Error(w, "Failed to load artist template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
