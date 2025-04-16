package handlers

import (
	"api/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"
	"context"

	"github.com/labstack/echo/v4"
	"cloud.google.com/go/firestore"
)

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func RefreshSpotifyToken(refreshToken, clientID, clientSecret string) (*SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	if tokenResp.AccessToken == "" {
		return nil, errors.New("invalid response from Spotify")
	}

	return &tokenResp, nil
}

func GetSpotifyTokenHandler(client *firestore.Client, clientID, clientSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		// search for the user in Firestore
		docRef := client.Collection("users").Doc(uid)
		docSnapshot, err := docRef.Get(context.Background())
		if err != nil {
			log.Println("Error getting user:", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		// extract the token from the document
		var token models.SpotifyToken
		if err := docSnapshot.DataTo(&token); err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Token not found"})
		}

		log.Println("Token found:", token)

		// check if the token is expired
		if time.Now().After(token.ExpiresAt) {
			newToken, err := RefreshSpotifyToken(token.RefreshToken, clientID, clientSecret)
			if err != nil {
				log.Println("Error refreshing token:", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to refresh token"})
			}

			// update token in Firestore
			token.AccessToken = newToken.AccessToken
			token.ExpiresAt = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)
			token.Scope = newToken.Scope
			token.TokenType = newToken.TokenType

			if newToken.RefreshToken != "" {
				token.RefreshToken = newToken.RefreshToken
			}

			// save the updated token in Firestore
			_, err = docRef.Set(context.Background(), token)
			if err != nil {
				log.Println("Error updating token in Firestore:", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update token"})
			}
		}

		// send the token back to the client
		return c.JSON(http.StatusOK, map[string]interface{}{
			"access_token": token.AccessToken,
			"expires_at":   token.ExpiresAt,
			"token_type":   token.TokenType,
			"scope":        token.Scope,
		})
	}
}
