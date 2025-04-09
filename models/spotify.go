package models

import "time"

type SpotifyToken struct {
	UserID       string    `gorm:"primaryKey"`
	AccessToken  string    `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	Scope        string
	TokenType    string    `gorm:"default:Bearer"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

