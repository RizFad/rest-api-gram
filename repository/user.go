package repository

import (
	"context"

	"mygram/infrastructure"
	"mygram/model"

	"gorm.io/gorm"
)

type UserQuery interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersByID(ctx context.Context, id uint64) (model.User, error)
	GetUsersByUsername(ctx context.Context, email string) (model.User, error)

	DeleteUsersByID(ctx context.Context, id uint64) error
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	UpdateUserByID(ctx context.Context, id uint64, user model.User) (model.User, error)
}

type UserCommand interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}

type userQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewUserQuery(db infrastructure.GormPostgres) UserQuery {
	return &userQueryImpl{db: db}
}

func (u *userQueryImpl) GetUsersByUsername(ctx context.Context, email string) (model.User, error) {
	db := u.db.GetConnection()
	user := model.User{}
	if err := db.WithContext(ctx).Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		// If user not found, return nil error
		if err == gorm.ErrRecordNotFound {
			return model.User{}, nil
		}
		return model.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	db := u.db.GetConnection()
	users := []model.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userQueryImpl) GetUsersByID(ctx context.Context, id uint64) (model.User, error) {
	db := u.db.GetConnection()
	users := model.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		Find(&users).Error; err != nil {
		// if user not found, return nil error
		if err == gorm.ErrRecordNotFound {
			return model.User{}, nil
		}

		return model.User{}, err
	}
	return users, nil
}

func (u *userQueryImpl) DeleteUsersByID(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Delete(&model.User{ID: id}).
		Error; err != nil {
		return err
	}
	return nil
}

func (u *userQueryImpl) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Save(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userQueryImpl) UpdateUserByID(ctx context.Context, id uint64, user model.User) (model.User, error) {
	db := u.db.GetConnection()
	// Update user by ID
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		Updates(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil

}
