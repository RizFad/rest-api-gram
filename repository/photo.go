package repository

import (
	"context"
	"mygram/infrastructure"
	"mygram/model"

	"gorm.io/gorm"
)

type PhotosQuery interface {
	GetAllPhotos(ctx context.Context) ([]model.Photo, error)
	UpdatePhoto(ctx context.Context, currentPhoto, newPhoto *model.Photo) (*model.Photo, error)
	DeletePhoto(ctx context.Context, photo *model.Photo) error
	FindPhotoByID(ctx context.Context, photoId int) (*model.Photo, error)
	CreatePhoto(ctx context.Context, photo *model.Photo) (*model.Photo, error)
}

type PhotoCommand interface {
	CreatePhoto(ctx context.Context, photo *model.Photo) (*model.Photo, error)
}

type photoQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewPhotoQuery(db infrastructure.GormPostgres) PhotosQuery {
	return &photoQueryImpl{db: db}
}

func (p *photoQueryImpl) CreatePhoto(ctx context.Context, photo *model.Photo) (*model.Photo, error) {
	err := p.db.GetConnection().Create(photo).Error
	if err != nil {
		return nil, err
	}

	return photo, err
}

func (p *photoQueryImpl) GetAllPhotos(ctx context.Context) ([]model.Photo, error) {
	var photos []model.Photo

	db := p.db.GetConnection()

	err :=
		db.WithContext(ctx).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Email", "Username")
		}).Find(&photos).Error

	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (p *photoQueryImpl) UpdatePhoto(ctx context.Context, currentPhoto, newPhoto *model.Photo) (*model.Photo, error) {
	db := p.db.GetConnection()
	// Update photo by ID
	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", currentPhoto.ID).
		Updates(&newPhoto).Error; err != nil {
		return &model.Photo{}, err
	}
	return newPhoto, nil
}

func (p *photoQueryImpl) DeletePhoto(ctx context.Context, photo *model.Photo) error {
	db := p.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("photos").
		Delete(&model.Photo{ID: photo.ID}).
		Error; err != nil {
		return err
	}
	return nil
}

func (p *photoQueryImpl) FindPhotoByID(ctx context.Context, photoId int) (*model.Photo, error) {
	db := p.db.GetConnection()
	photo := &model.Photo{}

	if err := db.
		WithContext(ctx).
		Table("photos").
		Where("id = ?", photoId).
		First(photo).Error; err != nil {

		// if photo not found, return nil error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}
	return photo, nil
}
