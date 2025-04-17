package models

import (
	"time"
)

// ShareLink 分享链接模型
type ShareLink struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	NoteID    string    `json:"note_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Note 笔记模型
type Note struct {
	ID         string      `json:"id" gorm:"primaryKey"`
	UserID     uint        `json:"user_id" gorm:"not null"`
	CategoryID string      `json:"category_id"`
	Title      string      `json:"title" gorm:"not null"`
	Content    string      `json:"content" gorm:"type:text;not null"`
	Category   *Category   `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags       []Tag       `json:"tags,omitempty" gorm:"many2many:note_tags;"`
	ShareLinks []ShareLink `json:"share_links,omitempty" gorm:"foreignKey:NoteID"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// NoteTag 笔记标签关联表
type NoteTag struct {
	NoteID string `gorm:"primaryKey"`
	TagID  string `gorm:"primaryKey"`
}
