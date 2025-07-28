package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	DiscordID string `gorm:"unique"`
}
