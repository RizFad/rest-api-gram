package repository

import (
	"context"
	"errors"
	"fmt"
	"mygram/infrastructure"
	"mygram/model"

	"gorm.io/gorm"
)

type SocialMediasQuery interface {
	CreateSocialMedia(ctx context.Context, socialMedia *model.SocialMedias) (*model.SocialMedias, error)
	GetAllSocialMedia(ctx context.Context) ([]model.SocialMedias, error)
	UpdateSocialMedia(ctx context.Context, currentsocialMedia, newsocialMedia *model.SocialMedias) (*model.SocialMedias, error)
	DeleteSocialMedia(ctx context.Context, socialMedia *model.SocialMedias) error
	FindSocialMediaByID(ctx context.Context, id int) (*model.SocialMedias, error)
}

type SocialMediasCommand interface {
	CreateSocialMedia(ctx context.Context, socialMedia *model.SocialMedias) (*model.SocialMedias, error)
}

type socialmediasQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewSocialMediasQuery(db infrastructure.GormPostgres) SocialMediasQuery {
	return &socialmediasQueryImpl{db: db}
}

func (sm *socialmediasQueryImpl) CreateSocialMedia(ctx context.Context, socialMedia *model.SocialMedias) (*model.SocialMedias, error) {
	err := sm.db.GetConnection().Create(socialMedia).Error
	if err != nil {
		return nil, err
	}

	return socialMedia, err
}

func (sm *socialmediasQueryImpl) GetAllSocialMedia(ctx context.Context) ([]model.SocialMedias, error) {
	var socialMedia []model.SocialMedias

	db := sm.db.GetConnection()

	err :=
		db.WithContext(ctx).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Email", "Username")
		}).Find(&socialMedia).Error

	if err != nil {
		return nil, err
	}

	return socialMedia, nil
}

func (sm *socialmediasQueryImpl) UpdateSocialMedia(ctx context.Context, currentsocialMedia, newsocialMedia *model.SocialMedias) (*model.SocialMedias, error) {
	db := sm.db.GetConnection()

	err := db.WithContext(ctx).Model(&currentsocialMedia).Updates(&newsocialMedia).Find(&currentsocialMedia).Error
	if err != nil {
		return nil, err
	}
	return currentsocialMedia, nil
}

func (sm *socialmediasQueryImpl) DeleteSocialMedia(ctx context.Context, socialMedia *model.SocialMedias) error {
	db := sm.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("social_medias").
		Delete(&model.SocialMedias{ID: socialMedia.ID}).
		Error; err != nil {
		return err
	}
	return nil
}

func (sm *socialmediasQueryImpl) FindSocialMediaByID(ctx context.Context, id int) (*model.SocialMedias, error) {
	db := sm.db.GetConnection()
	socialMedias := &model.SocialMedias{}

	err := db.WithContext(ctx).First(&socialMedias, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Social media with id %d not found.", id)
		}
		return nil, err
	}

	return socialMedias, nil
}
