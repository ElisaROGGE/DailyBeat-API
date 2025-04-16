package handlers

import (
	"api/models"
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
)

func SendFriendRequest(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get the user ID from the context
		fromUID, ok := c.Get("uid").(string)
		if !ok || fromUID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Utilisateur non authentifié"})
		}

		// get the payload
		var payload struct {
			ToUID string `json:"toUID"`
		}
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payload invalide"})
		}

		if payload.ToUID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "toUID manquant"})
		}

		// can't send a friend request to yourself
		if fromUID == payload.ToUID {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "On ne peut pas s’ajouter soi-même"})
		}

		ctx := context.Background()
		_, _, err := client.Collection("friendships").Add(ctx, models.Friendship{
			From:      fromUID,
			To:        payload.ToUID,
			Status:    "pending",
			CreatedAt: time.Now(),
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur Firestore"})
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "Demande d’ami envoyée"})
	}
}

