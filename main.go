package main

import (
	"api/config"
	"firebase.google.com/go/v4"
	"api/models"
	"api/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"context"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur lors du chargement de .env")
	}

	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Erreur d'initialisation de Firebase: %v", err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Erreur d'initialisation de Firestore: %v", err)
	}
	defer client.Close()

	db := config.ConnectDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	err = db.AutoMigrate(&models.User{}, &models.SpotifyToken{})
	if err != nil {
		log.Fatal("Erreur lors de la migration des modèles : ", err)
	}

	e := echo.New()

	routes.RegisterRoutes(e, app, db, client)

	e.Logger.Fatal(e.Start(":8080"))
}
