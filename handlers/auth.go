package handlers

import (
	"api/config"
	"api/models"
	"api/utils"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username       string `json:"username,omitempty"`
	Password       string `json:"password,omitempty"`
	SpotifyToken   string `json:"spotify_token,omitempty"`
	SpotifyRefresh string `json:"spotify_refresh,omitempty"`
	Country        string `json:"country"`
}

func CreateUser(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Email already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Database insertion failed"})
	}

	token, _ := utils.GenerateJWT(user.ID)

	return c.JSON(http.StatusCreated, map[string]string{"token": token})
}

func Login(c echo.Context) error {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var user models.User
	if err := config.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Token generation failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// Redirect to Spotify OAuth
func SpotifyLogin(c echo.Context) error {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	scope := "user-read-email user-read-private"

	authURL := fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s",
		clientID, url.QueryEscape(redirectURI), url.QueryEscape(scope),
	)

	return c.Redirect(http.StatusFound, authURL)
}

func LinkSpotifyAccount(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	fmt.Println("Authorization Header: ", token)
	token = strings.TrimPrefix(token, "Bearer ")
	userID, err := utils.ParseJWT(token)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get values"})
	}

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Authorization code is missing"})
	}

	accessToken, refreshToken, err := utils.ExchangeSpotifyToken(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get token"})
	}

	spotifyUser, err := utils.GetSpotifyUserID(accessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch Spotify user"})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	user.SpotifyID = spotifyUser
	user.SpotifyToken = accessToken
	user.SpotifyRefresh = refreshToken
	config.DB.Save(&user)

	return c.JSON(http.StatusOK, map[string]string{"message": "Spotify account linked successfully"})
}

func RefreshSpotify(c echo.Context) error {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	newToken, err := utils.RefreshSpotifyToken(request.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to refresh token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"spotify_token": newToken})
}
