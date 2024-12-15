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

func CreateUser(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Data connection failed"})
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Database insertion failed"})
	}

	return c.JSON(http.StatusCreated, user)
}
