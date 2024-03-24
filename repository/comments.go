package repository

import (
	"context"
	"errors"
	"fmt"
	"mygram/infrastructure"
	"mygram/model"

	"gorm.io/gorm"
)

type CommentsQuery interface {
	CreateComment(ctx context.Context, comment *model.Comments) (*model.Comments, error)
	GetAllComment(ctx context.Context) ([]model.Comments, error)
	UpdateComment(ctx context.Context, currentComment, newComment *model.Comments) (*model.Comments, error)
	DeleteComment(ctx context.Context, comment *model.Comments) error
	FindCommentByID(ctx context.Context, id int) (*model.Comments, error)
}

type CommentsCommand interface {
	CreateComment(ctx context.Context, comment *model.Comments) (*model.Comments, error)
}

type commentsQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewCommentsQuery(db infrastructure.GormPostgres) CommentsQuery {
	return &commentsQueryImpl{db: db}
}

func (c *commentsQueryImpl) CreateComment(ctx context.Context, comment *model.Comments) (*model.Comments, error) {
	err := c.db.GetConnection().Create(comment).Error
	if err != nil {
		return nil, err
	}

	return comment, err
}

func (c *commentsQueryImpl) GetAllComment(ctx context.Context) ([]model.Comments, error) {
	var comments []model.Comments

	db := c.db.GetConnection()

	err :=
		db.WithContext(ctx).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Email", "Username")
		}).Preload("Photo").Find(&comments).Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *commentsQueryImpl) UpdateComment(ctx context.Context, currentComment, newComment *model.Comments) (*model.Comments, error) {
	db := c.db.GetConnection()

	err := db.WithContext(ctx).Preload("Photo").Model(&currentComment).Updates(&newComment).Find(&currentComment).Error
	if err != nil {
		return nil, err
	}
	return currentComment, nil
}

func (c *commentsQueryImpl) DeleteComment(ctx context.Context, comment *model.Comments) error {
	db := c.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("comments").
		Delete(&model.Comments{ID: comment.ID}).
		Error; err != nil {
		return err
	}
	return nil
}

func (c *commentsQueryImpl) FindCommentByID(ctx context.Context, id int) (*model.Comments, error) {
	db := c.db.GetConnection()
	comment := &model.Comments{}

	err := db.WithContext(ctx).First(&comment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Comment with id %d not found.", id)
		}
		return nil, err
	}

	return comment, nil
}
