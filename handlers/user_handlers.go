package handlers

import (
    "api/config"
    "api/models"
    "net/http"

    "github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
    var users []models.User
    config.DB.Find(&users)
    return c.JSON(http.StatusOK, users)
}
