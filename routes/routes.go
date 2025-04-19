package routes

import (
	"api/handlers"
	"api/middleware"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, app *firebase.App, db *gorm.DB, client *firestore.Client) {
	e.GET("/users", handlers.GetUsers)
	e.POST("/users", handlers.CreateUser)

	auth := e.Group("/auth")
	auth.Use(middleware.FirebaseAuthMiddleware(app))
	auth.GET("/spotify-token", handlers.GetSpotifyTokenHandler(client, os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET")))
	auth.POST("/friends/request", handlers.SendFriendRequest(client))
	auth.POST("/friends/accept", handlers.AcceptFriendRequest(client))
	auth.POST("/music-of-the-day", handlers.SetMusicOfTheDay(client))

}

