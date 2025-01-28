package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var (
	clientID     string
	clientSecret string
)

func InitAPI() {
	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
}

// Obtenir le token Spotify
func GetSpotifyToken() (string, error) {
	authURL := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	accessToken, ok := response["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to retrieve access token")
	}

	return accessToken, nil
}

// Chercher un artiste via Spotify
func SearchArtist(artistName string, token string) (string, error) {
	searchURL := "https://api.spotify.com/v1/search"

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", err
	}

	// Ajout du paramètre de recherche
	q := req.URL.Query()
	q.Add("q", artistName)
	q.Add("type", "artist")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Extraire l'ID de l'artiste
	artists := response["artists"].(map[string]interface{})["items"].([]interface{})
	if len(artists) > 0 {
		artist := artists[0].(map[string]interface{})
		return artist["id"].(string), nil
	}

	return "", fmt.Errorf("artist not found")
}
