package routes

import (
    "api/handlers"

    "github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
    e.GET("/users", handlers.GetUsers)
    e.POST("/users", handlers.CreateUser)
    e.GET("/spotify/callback", handlers.HandleSpotifyCallback) 
    e.POST("/auth/register", handlers.CreateUser)
    e.POST("/auth/login", handlers.Login)
    e.GET("/auth/spotify/link", handlers.LinkSpotifyAccount)
}
