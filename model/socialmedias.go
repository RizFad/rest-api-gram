package model

import (
	"errors"
	"time"
)

type SocialMedias struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"notNull"`
	URL       string    `json:"url" gorm:"notNull"`
	UserID    int       `json:"user_id" gorm:"notNull"`
	User      User      `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SocialMediaUserGet struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type SocialMediaGet struct {
	ID        int                `json:"id" gorm:"primaryKey"`
	Name      string             `json:"name" gorm:"notNull"`
	URL       string             `json:"url" gorm:"notNull"`
	UserID    int                `json:"user_id" gorm:"notNull"`
	User      SocialMediaUserGet `json:"user"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type SocialMediaUpdate struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"notNull"`
	URL       string    `json:"url" gorm:"notNull"`
	UserID    int       `json:"user_id" gorm:"notNull"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SocialMediaCreate struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required,url"`
}

type UpdateSocialMedia struct {
	Name string `json:"name"`
	URL  string `json:"url" validate:"url"`
}

func (sm SocialMediaCreate) SocialMediaValidate() error {
	if sm.Name == "" {
		return errors.New("invalid name cause is required")
	}
	if sm.URL == "" {
		return errors.New("invalid social media url cause is required")
	}
	return nil
}
