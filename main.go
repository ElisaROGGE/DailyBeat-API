package main

import (
	"api/config"
	"api/firebase"
	"api/models"
	"api/routes"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur lors du chargement de .env")
	}
	firebase.InitFirebase()
	app := firebase.FirebaseApp


	db := config.ConnectDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&models.User{})
	err = db.AutoMigrate(&models.SpotifyToken{})
	if err != nil {
		log.Fatal("Erreur lors de la migration des modèles : ", err)
	}

	e := echo.New()

	routes.RegisterRoutes(e, app, db)

	e.Logger.Fatal(e.Start(":8080"))
}
