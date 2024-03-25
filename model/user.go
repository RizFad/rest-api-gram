package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	DoB       time.Time      `json:"age" gorm:"column:dob"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}

type DefaultColumn struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

type UserMediaSocial struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type UserSignUp struct {
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password" binding:"required"`
	Email    string    `json:"email"`
	DoB      time.Time `json:"age"`
}

type UserSignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u UserSignUp) Validate() error {
	if u.Username == "" {
		return errors.New("invalid username")
	}
	if len(u.Password) < 6 {
		return errors.New("invalid password: length must be at least 6 characters")
	}

	dobString := u.DoB.Format(time.RFC3339)
	dob, err := time.Parse(time.RFC3339, dobString)
	if err != nil {
		return errors.New("invalid age format: must be in format (e.g., 2020-01-01T00:00:00Z)")
	}

	today := time.Now()
	age := today.Year() - dob.Year()

	if age < 8 || age > 150 {
		return errors.New("invalid age: must be between 8 and 150 years old")
	}

	return nil
}

func (u UserSignIn) Authenticate(passwordHash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(u.Password)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}
