package main

import (
	"log"
	"api/config"
	"api/routes"
	"api/models"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur lors du chargement de .env")
	}

	db := config.ConnectDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Erreur lors de la migration des modèles : ", err)
	}

	e := echo.New()

	routes.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
