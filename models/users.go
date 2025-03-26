package models

import "time"

type User struct {
	ID              uint   `gorm:"primaryKey"`
	Username        string `json:"username,omitempty"`
	Email        string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	SpotifyID       string `json:"spotify_id,omitempty" gorm:"unique"`
	SpotifyToken    string `json:"spotify_token,omitempty"`
	SpotifyRefresh  string `json:"spotify_refresh,omitempty"`
	Country         string `json:"country"`
	CreatedAt       time.Time `json:"created_at"`
}
