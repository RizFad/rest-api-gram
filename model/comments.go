package model

import (
	"time"
)

type Comments struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Message   string    `json:"message" gorm:"notNull"`
	PhotoID   int       `json:"photo_id" gorm:"notNull"`
	UserID    int       `json:"user_id" gorm:"notNull"`
	User      User      `json:"-"`
	Photo     Photo     `json:"-"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type CommentUser struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type CommentPhoto struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
	UserID  int    `json:"user_id"`
}

type CommentGetAll struct {
	ID        int          `json:"id" gorm:"primaryKey"`
	Message   string       `json:"message" gorm:"notNull"`
	PhotoID   int          `json:"photo_id" gorm:"notNull"`
	UserID    int          `json:"user_id" gorm:"notNull"`
	User      CommentUser  `json:"user"`
	Photo     CommentPhoto `json:"photo"`
	CreatedAt time.Time    `json:"create_at"`
	UpdatedAt time.Time    `json:"update_at"`
}

type CommentUpdate struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	URL       string    `json:"url"`
	UserID    int       `json:"user_id"`
	UpdatedAt time.Time `json:"update_at"`
}

type CreateComment struct {
	Message string `json:"message" validate:"required"`
	PhotoID int    `json:"photo_id" validate:"required"`
}

type UpdateComment struct {
	Message string `json:"message"`
}
