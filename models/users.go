package models

import "time"

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username     string `json:"username"`
    Token    string `json:"token" gorm:"unique"`
    Country string `json:"country"`
	CreatedAt time.Time `json:"created_at"`
}
