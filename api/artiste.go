package api

import (
	"html/template"
	"net/http"
)

type ArtistPageData struct {
	Name           string
	SpotifyArtistID string
	CreationDate   int
	Members        []string
	Image          string
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
		Name:           artistData.Name,
		SpotifyArtistID: spotifyID,
		CreationDate:   artistData.CreationDate,
		Members:        artistData.Members,
		Image:          artistData.Image,
	}

	tmpl, err := template.ParseFiles("./html/artiste.html")
	if err != nil {
		http.Error(w, "Failed to load artist template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
