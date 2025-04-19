package handlers

import (
	"cloud.google.com/go/firestore"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type MusicOfTheDayPayload struct {
	TrackID   string `json:"trackID"`
	TrackName string `json:"trackName"`
	Artist    string `json:"artist"`
}

func SetMusicOfTheDay(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		userUID, ok := c.Get("uid").(string)
		if !ok || userUID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Utilisateur non authentifié"})
		}

		var payload MusicOfTheDayPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payload invalide"})
		}

		if payload.TrackID == "" || payload.TrackName == "" || payload.Artist == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Champs manquants"})
		}

		// add the music of the day to Firestore
		ctx := context.Background()
		_, err := client.Collection("music_of_the_day").Doc(userUID).Set(ctx, map[string]interface{}{
			"trackID":   payload.TrackID,
			"trackName": payload.TrackName,
			"artist":    payload.Artist,
			"addedAt":   time.Now(),
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur lors de l'ajout de la musique"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Musique du jour ajoutée avec succès"})
	}
}
