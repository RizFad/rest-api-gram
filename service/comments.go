package service

import (
	"context"
	"fmt"
	"mygram/model"
	"mygram/repository"
)

type CommentsService interface {
	GetAllComment(ctx context.Context) ([]model.CommentGetAll, error)
	UpdateComment(ctx context.Context, data model.UpdateComment, commentID, userID int) (*model.CommentUpdate, error)
	CreateComment(ctx context.Context, data model.CreateComment, userId int) (*model.Comments, error)
	DeleteComment(ctx context.Context, commentID int, userID int) error
}

type commentsServiceImpl struct {
	repo repository.CommentsQuery
}

func NewCommentsService(repo repository.CommentsQuery) CommentsService {
	return &commentsServiceImpl{repo: repo}
}

func (c *commentsServiceImpl) GetAllComment(ctx context.Context) ([]model.CommentGetAll, error) {
	comments, err := c.repo.GetAllComment(ctx)
	if err != nil {
		return nil, err
	}

	var dataComment []model.CommentGetAll
	for _, comment := range comments {
		newComment := model.CommentGetAll{
			ID:        comment.ID,
			Message:   comment.Message,
			PhotoID:   comment.PhotoID,
			UserID:    comment.UserID,
			CreatedAt: comment.CreatedAt,
			User: model.CommentUser{
				ID:       comment.User.ID,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo: model.CommentPhoto{
				ID:      comment.Photo.ID,
				Title:   comment.Photo.Title,
				Caption: comment.Photo.Caption,
				URL:     comment.Photo.URL,
				UserID:  comment.Photo.UserID,
			},
		}
		dataComment = append(dataComment, newComment)
	}
	return dataComment, nil
}

func (c *commentsServiceImpl) UpdateComment(ctx context.Context, data model.UpdateComment, commentID, userID int) (*model.CommentUpdate, error) {
	currentComment, err := c.repo.FindCommentByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if currentComment.UserID != userID {
		return nil, fmt.Errorf("comment with id %d is not a comment owned by user with id %d.", commentID, userID)
	}

	newComment := &model.Comments{Message: data.Message}

	updatedPhoto, err := c.repo.UpdateComment(ctx, currentComment, newComment)
	if err != nil {
		return nil, err
	}

	dataComment := &model.CommentUpdate{
		ID:        updatedPhoto.ID,
		Title:     updatedPhoto.Photo.Title,
		Caption:   updatedPhoto.Photo.Caption,
		URL:       updatedPhoto.Photo.URL,
		UserID:    updatedPhoto.UserID,
		UpdatedAt: updatedPhoto.UpdatedAt,
	}

	return dataComment, nil
}

func (c *commentsServiceImpl) DeleteComment(ctx context.Context, commentID int, userID int) error {
	comment, err := c.repo.FindCommentByID(ctx, commentID)
	if err != nil {
		return err

	}

	if comment.UserID != userID {
		return fmt.Errorf("Comment with id %d is not a comment owned by user with id %d.", commentID, userID)
	}

	err = c.repo.DeleteComment(ctx, comment)
	if err != nil {
		return fmt.Errorf("Error deleting comment: %v", err)
	}

	return err
}

func (c *commentsServiceImpl) CreateComment(ctx context.Context, data model.CreateComment, userId int) (*model.Comments, error) {
	comment := &model.Comments{
		Message: data.Message,
		PhotoID: data.PhotoID,
		UserID:  userId,
	}

	dataComment, err := c.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return dataComment, nil
}
