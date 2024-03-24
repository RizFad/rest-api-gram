package model

import (
	"time"
)

type Photo struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"notNull"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"url" gorm:"notNull"`
	UserID    int       `json:"user_id" gorm:"notNull"`
	User      User      `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PhotoUserGet struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type CreatePhoto struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" validate:"required,url"`
}

type UpdatePhoto struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
}

type PhotoGet struct {
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	Caption   string       `json:"caption"`
	PhotoURL  string       `json:"photo_url"`
	UserID    int          `json:"user_id"`
	User      PhotoUserGet `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type PhotoUpdate struct {
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UserID    int       `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}
