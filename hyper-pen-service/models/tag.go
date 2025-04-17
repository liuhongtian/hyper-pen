package models

import (
	"time"
)

type Tag struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:255;not null"`
	Color     string    `json:"color" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"size:36;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Notes     []Note    `json:"notes" gorm:"many2many:note_tags;"`
}
