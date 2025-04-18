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

		docRef := client.Collection("users").Doc(payload.ToUID)
		docSnap, err := docRef.Get(ctx)
		if err != nil || !docSnap.Exists() {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Utilisateur cible introuvable"})
		}

		_, _, err = client.Collection("friendships").Add(ctx, models.Friendship{
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

func AcceptFriendRequest(client *firestore.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		toUID, ok := c.Get("uid").(string)
		if !ok || toUID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Utilisateur non authentifié"})
		}

		var payload struct {
			FromUID string `json:"fromUID"`
		}
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payload invalide"})
		}

		if payload.FromUID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "fromUID manquant"})
		}

		ctx := context.Background()

		query := client.Collection("friendships").
			Where("from", "==", payload.FromUID).
			Where("to", "==", toUID).
			Where("status", "==", "pending")

		docs, err := query.Documents(ctx).GetAll()
		if err != nil || len(docs) == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Demande d’ami non trouvée"})
		}

		_, err = docs[0].Ref.Update(ctx, []firestore.Update{
			{Path: "status", Value: "accepted"},
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Échec de l’acceptation"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Demande d’ami acceptée"})
	}
}

