package models

import "time"

type FriendRequestPayload struct {
	FromUID string `json:"fromUID"`
	ToUID   string `json:"toUID"`
}

type Friendship struct {
	From      string    `firestore:"from"`
	To        string    `firestore:"to"`
	Status    string    `firestore:"status"`
	CreatedAt time.Time `firestore:"createdAt"`
}
