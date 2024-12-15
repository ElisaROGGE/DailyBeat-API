package routes

import (
    "api/handlers"

    "github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
    e.GET("/users", handlers.GetUsers)
    e.POST("/users", handlers.CreateUser)
}
