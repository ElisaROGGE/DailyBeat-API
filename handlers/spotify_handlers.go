package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"
	"github.com/labstack/echo/v4"
)

// SpotifyTokenResponse structure
type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func HandleSpotifyCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.String(http.StatusBadRequest, "Code manquant")
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))
	data.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("SPOTIFY_CLIENT_SECRET"))

	// send request to Spotify
	resp, err := http.PostForm("https://accounts.spotify.com/api/token", data)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Erreur requête Spotify")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.String(http.StatusUnauthorized, "Impossible d'obtenir le token")
	}

	// read JSON response
	var tokenResponse SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return c.String(http.StatusInternalServerError, "Erreur décodage JSON")
	}

	// define HttpOnly cookie to store token
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     "spotify_access_token",
		Value:    tokenResponse.AccessToken,
		Path:     "/",
		HttpOnly: true, 
		Secure:   true, 
		Expires:  time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
	})

	return c.Redirect(http.StatusSeeOther, "/")
}
