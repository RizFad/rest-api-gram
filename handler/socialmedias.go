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

type SocialMediasHandler interface {
	CreateSocialMedia(ctx *gin.Context)
	GetAllSocialMedia(ctx *gin.Context)
	UpdateSocialMedia(ctx *gin.Context)
	DeleteSocialMedia(ctx *gin.Context)
}

type socialmediasHandlerImpl struct {
	svc service.SocialMediasService
}

func NewSocialMediasHandler(svc service.SocialMediasService) SocialMediasHandler {
	return &socialmediasHandlerImpl{
		svc: svc,
	}
}

func (sm *socialmediasHandlerImpl) UpdateSocialMedia(ctx *gin.Context) {
	var data model.UpdateSocialMedia

	socialmediaID, err := strconv.Atoi(ctx.Param("socialMediaId"))
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

	// Call service to update socialmedia
	updatedSocialMedias, err := sm.svc.UpdateSocialMedia(ctx, data, socialmediaID, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// Respond with updated socialmedia details
	ctx.JSON(http.StatusOK, updatedSocialMedias)
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
func (sm *socialmediasHandlerImpl) GetAllSocialMedia(ctx *gin.Context) {
	socialmedias, err := sm.svc.GetAllSocialMedia(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, socialmedias)
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
func (sm *socialmediasHandlerImpl) CreateSocialMedia(ctx *gin.Context) {
	socialmediaCreate := model.SocialMediaCreate{}

	err := ctx.ShouldBindJSON(&socialmediaCreate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid request body"})
		return
	}

	errs := socialmediaCreate.SocialMediaValidate()
	if errs != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid request body"})
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

	// Call service to create socialmedia
	socialmedia, err := sm.svc.CreateSocialMedia(ctx, socialmediaCreate, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, socialmedia)
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
func (sm *socialmediasHandlerImpl) DeleteSocialMedia(ctx *gin.Context) {
	// Get social media ID
	socialmediaID, err := strconv.Atoi(ctx.Param("socialMediaId"))
	if socialmediaID == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid social media ID"})
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

	// Check if the photo belongs to the user
	if socialmediaID != int(userId) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid request: social media does not belong to the user"})
		return
	}

	// Call service to delete photo
	err = sm.svc.DeleteSocialMedia(ctx, socialmediaID, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
