package utils

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/go-resty/resty/v2"
)

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func ExchangeSpotifyToken(code string) (string, string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  os.Getenv("SPOTIFY_REDIRECT_URI"),
			"client_id":     os.Getenv("SPOTIFY_CLIENT_ID"),
			"client_secret": os.Getenv("SPOTIFY_CLIENT_SECRET"),
		}).
		Post("https://accounts.spotify.com/api/token")

	if err != nil {
		return "", "", err
	}

	if resp.StatusCode() != 200 {
		return "", "", errors.New("failed to get Spotify token")
	}

	var tokenResponse SpotifyTokenResponse
	if err := json.Unmarshal(resp.Body(), &tokenResponse); err != nil {
		return "", "", err
	}

	if err != nil {
		return "", "",  err
	}

	return tokenResponse.AccessToken, tokenResponse.RefreshToken, err
}

func GetSpotifyUserID(accessToken string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get("https://api.spotify.com/v1/me")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", errors.New("failed to get Spotify user ID")
	}

	var userData struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(resp.Body(), &userData); err != nil {
		return "", err
	}

	return userData.ID, nil
}

func RefreshSpotifyToken(refreshToken string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
			"client_id":     os.Getenv("SPOTIFY_CLIENT_ID"),
			"client_secret": os.Getenv("SPOTIFY_CLIENT_SECRET"),
		}).
		Post("https://accounts.spotify.com/api/token")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", errors.New("failed to refresh Spotify token")
	}

	var tokenResponse SpotifyTokenResponse
	if err := json.Unmarshal(resp.Body(), &tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}
