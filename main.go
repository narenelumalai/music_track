package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// LastFmResponse represents the structure of the response from Last.fm API
type LastFmResponse struct {
	Tracks struct {
		Track []struct {
			Name       string `json:"name"`
			Duration   string `json:"duration"`
			Listeners  string `json:"listeners"`
			Mbid       string `json:"mbid"`
			Url        string `json:"url"`
			Streamable struct {
				Text      string `json:"#text"`
				Fulltrack string `json:"fulltrack"`
			} `json:"streamable"`
			Artist struct {
				Name string `json:"name"`
				Mbid string `json:"mbid"`
				Url  string `json:"url"`
			} `json:"artist"`
			Image []struct {
				Text string `json:"#text"`
				Size string `json:"size"`
			} `json:"image"`
			Attr struct {
				Rank string `json:"rank"`
			} `json:"@attr"`
		} `json:"track"`
	} `json:"tracks"`
}

// HandleGetTopTracks handles GET requests to /toptracks endpoint
// @Summary Get top tracks for a location
// @Description Get the top tracks for a given location and country
// @Tags tracks
// @Produce json
// @Param location query string true "Location name"
// @Param country query string true "Country name"
// @Success 200 {object} LastFmResponse
// @Router /toptracks [get]
func HandleGetTopTracks(w http.ResponseWriter, r *http.Request) {
	// Get location and country from query parameters
	location := r.URL.Query().Get("location")
	country := r.URL.Query().Get("country")

	// Call Last.fm API
	url := fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=geo.gettoptracks&country=%s&location=%s&api_key=79d7732491b5ef4534fa5b90c5d27bc7&format=json", country, location)
	response, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Decode JSON response
	var lastFmResponse LastFmResponse
	err = json.NewDecoder(response.Body).Decode(&lastFmResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send JSON response back to the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lastFmResponse)
}

// @title Your API Title
// @version 1.0
// @description Your API description
// @termsOfService Your API terms of service
// @host localhost:8080
// @BasePath /v1
func main() {
	// Register HTTP handlers
	http.HandleFunc("/toptracks", HandleGetTopTracks)

	// Serve Swagger UI
	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
	))

	// Start HTTP server
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
