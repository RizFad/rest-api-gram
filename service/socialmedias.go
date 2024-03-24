package service

import (
	"context"
	"fmt"
	"mygram/model"
	"mygram/repository"
)

type SocialMediasService interface {
	CreateSocialMedia(ctx context.Context, data model.SocialMediaCreate, userID int) (*model.SocialMedias, error)
	GetAllSocialMedia(ctx context.Context) ([]model.SocialMediaGet, error)
	UpdateSocialMedia(ctx context.Context, data model.UpdateSocialMedia, smID int, userID int) (*model.SocialMediaUpdate, error)
	DeleteSocialMedia(ctx context.Context, smID, userID int) error
}

type socialmediasServiceImpl struct {
	repo repository.SocialMediasQuery
}

func NewSocialMediasService(repo repository.SocialMediasQuery) SocialMediasService {
	return &socialmediasServiceImpl{repo: repo}
}

func (sm *socialmediasServiceImpl) GetAllSocialMedia(ctx context.Context) ([]model.SocialMediaGet, error) {
	socialmedias, err := sm.repo.GetAllSocialMedia(ctx)
	if err != nil {
		return nil, err
	}
	respSocialMedias := parseSocialMediaGet(socialmedias)

	return respSocialMedias, nil
}

func (sm *socialmediasServiceImpl) UpdateSocialMedia(ctx context.Context, data model.UpdateSocialMedia, smID int, userID int) (*model.SocialMediaUpdate, error) {
	currentSocialMedia, err := sm.repo.FindSocialMediaByID(ctx, smID)
	if err != nil {
		return nil, err
	}

	if currentSocialMedia.UserID != userID {
		return nil, fmt.Errorf("Social Media with id %d is not a Social Media owned by user with id %d.", smID, userID)
	}

	newSocialMedia := &model.SocialMedias{
		Name: data.Name,
		URL:  data.URL,
	}

	updatedSocialMedia, err := sm.repo.UpdateSocialMedia(ctx, currentSocialMedia, newSocialMedia)
	if err != nil {
		return nil, err
	}

	respData := &model.SocialMediaUpdate{
		ID:        updatedSocialMedia.ID,
		Name:      updatedSocialMedia.Name,
		URL:       updatedSocialMedia.URL,
		UserID:    updatedSocialMedia.UserID,
		UpdatedAt: currentSocialMedia.UpdatedAt,
	}

	return respData, nil
}

func (sm *socialmediasServiceImpl) DeleteSocialMedia(ctx context.Context, smID, userID int) error {
	socialmedia, err := sm.repo.FindSocialMediaByID(ctx, smID)
	if err != nil {
		return err

	}

	if socialmedia.UserID != userID {
		return fmt.Errorf("Social Media  with id %d is not a social media owned by user with id %d.", socialmedia.UserID, userID)
	}

	err = sm.repo.DeleteSocialMedia(ctx, socialmedia)
	if err != nil {
		return fmt.Errorf("Error deleting social media: %v", err)
	}

	return err
}

func (sm *socialmediasServiceImpl) CreateSocialMedia(ctx context.Context, data model.SocialMediaCreate, userID int) (*model.SocialMedias, error) {
	socialmedia := &model.SocialMedias{
		UserID: userID,
		Name:   data.Name,
		URL:    data.URL,
	}

	resSocialMedias, err := sm.repo.CreateSocialMedia(ctx, socialmedia)
	if err != nil {
		return nil, err
	}

	return resSocialMedias, nil
}

func parseSocialMediaGet(socialMedias []model.SocialMedias) []model.SocialMediaGet {
	var parsedSocialMedia []model.SocialMediaGet
	for _, sm := range socialMedias {
		newSM := model.SocialMediaGet{
			ID:     sm.ID,
			Name:   sm.Name,
			URL:    sm.URL,
			UserID: sm.UserID,
			User: model.SocialMediaUserGet{
				Email:    sm.User.Email,
				Username: sm.User.Username,
			},
			CreatedAt: sm.CreatedAt,
			UpdatedAt: sm.UpdatedAt,
		}
		parsedSocialMedia = append(parsedSocialMedia, newSM)
	}
	return parsedSocialMedia
}
