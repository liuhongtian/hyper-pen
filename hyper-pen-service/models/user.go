package models

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Username    string    `json:"username" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	GithubID    string    `json:"github_id" gorm:"unique"`
	WechatID    string    `json:"wechat_id" gorm:"unique"`
	AvatarURL   string    `json:"avatar_url"`
	GitHubToken string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
