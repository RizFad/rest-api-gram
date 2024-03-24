package handler

import (
	"net/http"
	"strconv"

	"mygram/middleware"
	"mygram/model"
	"mygram/pkg"
	"mygram/service"

	"github.com/gin-gonic/gin"
)

type CommentsHandler interface {
	GetAllComment(ctx *gin.Context)
	UpdateComment(ctx *gin.Context)
	DeleteComment(ctx *gin.Context)
	CreateComment(ctx *gin.Context)
}

type commentHandlerImpl struct {
	svc service.CommentsService
}

func NewCommentHandler(svc service.CommentsService) CommentsHandler {
	return &commentHandlerImpl{
		svc: svc,
	}
}

func (c *commentHandlerImpl) UpdateComment(ctx *gin.Context) {
	var data model.UpdateComment

	commentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "ID must be a number"})
		return
	}

	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid request body"})
		return
	}

	// Get user from context
	user, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "user information not found in context"})
		return
	}

	userId, ok := user.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user ID in context"})
		return
	}

	// Call service to update photo
	updatedComment, err := c.svc.UpdateComment(ctx, data, commentID, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// Respond with updated photo details
	ctx.JSON(http.StatusOK, updatedComment)
}

// ShowUsers godoc
//
//	@Summary		Show users list
//	@Description	will fetch 3rd party server to get users data
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]model.User
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/users [get]
func (c *commentHandlerImpl) GetAllComment(ctx *gin.Context) {
	comments, err := c.svc.GetAllComment(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, comments)
}

// ShowUsersById godoc
//
//	@Summary		Show users detail
//	@Description	will fetch 3rd party server to get users data to get detail user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	model.User
//	@Failure		400	{object}	pkg.ErrorResponse
//	@Failure		404	{object}	pkg.ErrorResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/users/{id} [get]
func (c *commentHandlerImpl) CreateComment(ctx *gin.Context) {
	commentCreate := model.CreateComment{}

	err := ctx.ShouldBindJSON(&commentCreate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid request body"})
		return
	}

	// errs := commentCreate.CommentValidate()
	// if errs != nil {
	// 	ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid request body"})
	// 	return
	// }

	// Get user info from context
	user, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "user information not found in context"})
		return
	}

	userId, ok := user.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user ID in context"})
		return
	}

	// Call service to create comment
	comment, err := c.svc.CreateComment(ctx, commentCreate, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

// DeleteUsersById godoc
//
//		@Summary		Delete user by selected id
//		@Description	will delete user with given id from param
//		@Tags			users
//		@Accept			json
//		@Produce		json
//	 	@Param 			Authorization header string true "bearer token"
//		@Param			id	path		int	true	"User ID"
//		@Success		200	{object}	model.User
//		@Failure		400	{object}	pkg.ErrorResponse
//		@Failure		404	{object}	pkg.ErrorResponse
//		@Failure		500	{object}	pkg.ErrorResponse
//		@Router			/users/{id} [delete]
func (c *commentHandlerImpl) DeleteComment(ctx *gin.Context) {
	// Get comment ID
	commentID, err := strconv.Atoi(ctx.Param("id"))
	if commentID == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid photo ID"})
		return
	}

	// Get user info from context
	user, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "user information not found in context"})
		return
	}

	userId, ok := user.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user ID in context"})
		return
	}

	// Check if the comment belongs to the user
	if commentID != int(userId) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid request: comment does not belong to the user"})
		return
	}

	// Call service to delete comment
	err = c.svc.DeleteComment(ctx, commentID, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
