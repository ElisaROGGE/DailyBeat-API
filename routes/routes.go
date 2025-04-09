package routes

import (
	"api/handlers"
	"api/middleware"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, app *firebase.App, db *gorm.DB) {
	e.GET("/users", handlers.GetUsers)
	e.POST("/users", handlers.CreateUser)

	auth := e.Group("/auth")
	auth.Use(middleware.FirebaseAuthMiddleware(app))
	auth.GET("/spotify-token", handlers.GetSpotifyTokenHandler(db, os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET")))
}

