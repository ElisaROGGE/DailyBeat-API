package handlers

import (
	"api/models"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
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
func GetSpotifyTokenHandler(db *gorm.DB, clientID, clientSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.Param("uid")

		var token models.SpotifyToken
		if err := db.First(&token, "user_id = ?", uid).Error; err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Token not found"})
		}

		if time.Now().After(token.ExpiresAt) {
			newToken, err := RefreshSpotifyToken(token.RefreshToken, clientID, clientSecret)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
			}

			token.AccessToken = newToken.AccessToken
			token.ExpiresAt = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)
			token.Scope = newToken.Scope
			token.TokenType = newToken.TokenType

			if newToken.RefreshToken != "" {
				token.RefreshToken = newToken.RefreshToken
			}

			if err := db.Save(&token).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token"})
				
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"access_token": token.AccessToken,
			"expires_at":   token.ExpiresAt,
			"token_type":   token.TokenType,
			"scope":        token.Scope,
		})
	}
}


