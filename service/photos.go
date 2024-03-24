package service

import (
	"context"
	"fmt"
	"mygram/model"
	"mygram/repository"
)

type PhotosService interface {
	GetAllPhotos(ctx context.Context) ([]model.Photo, error)
	UpdatePhoto(ctx context.Context, currentPhoto, newPhoto *model.Photo) (*model.Photo, error)
	DeletePhoto(ctx context.Context, photo *model.Photo) error
	FindPhotoByID(ctx context.Context, photoId int) (*model.Photo, error)
	CreatePhoto(ctx context.Context, photo *model.Photo) (*model.Photo, error)
}

type photosServiceImpl struct {
	repo repository.PhotosQuery
}

func NewPhotosService(repo repository.PhotosQuery) PhotosService {
	return &photosServiceImpl{repo: repo}
}

func (p *photosServiceImpl) GetAllPhotos(ctx context.Context) ([]model.PhotoGet, error) {
	photos, err := p.repo.GetAllPhotos(ctx)
	if err != nil {
		return nil, err
	}
	respPhotos := parseGetAllPhotos(photos)

	return respPhotos, nil
}

func (p *photosServiceImpl) UpdatePhoto(ctx context.Context, req model.UpdatePhoto, photoId, userID int) (*model.PhotoUpdate, error) {
	currentPhoto, err := p.repo.FindPhotoByID(ctx, photoId)
	if err != nil {
		return nil, err
	}

	if currentPhoto.UserID != userID {
		return nil, fmt.Errorf("Photo with id %d is not a photo owned by user with id %d.", photoId, userID)
	}

	newPhoto := &model.Photo{
		PhotoURL: req.PhotoURL,
		Caption:  req.Caption,
		Title:    req.Caption,
	}

	updatedPhoto, err := p.repo.UpdatePhoto(ctx, currentPhoto, newPhoto)
	if err != nil {
		return nil, err
	}

	responsePhoto := parseUpdatePhoto(updatedPhoto)

	return responsePhoto, nil
}

func (p *photosServiceImpl) DeletePhoto(photoID int, userID int) error {
	photo, err := p.repo.FindPhotoByID(photoID)
	if err != nil {
		return err

	}

	if photo.UserID != userID {
		return fmt.Errorf("Photo with id %d is not a photo owned by user with id %d.", photoID, userID)
	}

	err = p.repo.DeletePhoto(photo)
	if err != nil {
		return fmt.Errorf("Error deleting photo: %v", err)
	}

	return err
}

func (p *photosServiceImpl) FindPhotoByID(ctx context.Context, id int) (model.Photo, error) {
	photo, err := p.repo.FindPhotoByID(ctx, id)
	if err != nil {
		return model.Photo{}, err
	}
	return *photo, err
}

func (p *photosServiceImpl) CreatePhoto(req model.CreatePhoto, userId int) (model.Photo, error) {
	photo := &model.Photo{
		Title:    req.Title,
		Caption:  req.Caption,
		PhotoURL: req.PhotoURL,
		UserID:   userId,
	}

	resPhoto, err := p.repo.CreatePhoto(photo)
	if err != nil {
		return nil, err
	}

	return resPhoto, nil
}

func parseGetAllPhotos(photos []model.Photo) []model.PhotoGet {
	var parsedPhotos []model.PhotoGet
	for _, photo := range photos {
		newPhoto := model.PhotoGet{
			ID:       photo.ID,
			Title:    photo.Title,
			Caption:  photo.Caption,
			PhotoURL: photo.PhotoURL,
			UserID:   photo.UserID,
			User: model.PhotoUserGet{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
		}
		parsedPhotos = append(parsedPhotos, newPhoto)
	}
	return parsedPhotos
}

func parseUpdatePhoto(photo *model.Photo) *model.PhotoUpdate {
	updatedPhoto := &model.PhotoUpdate{
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoURL:  photo.PhotoURL,
		UserID:    photo.UserID,
		UpdatedAt: photo.UpdatedAt,
	}

	return updatedPhoto
}
