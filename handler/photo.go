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

type PhotoHandler interface {
	GetAllPhotos(ctx *gin.Context)
	UpdatePhoto(ctx *gin.Context)
	DeletePhoto(ctx *gin.Context)
	CreatePhoto(ctx *gin.Context)
}

type photoHandlerImpl struct {
	svc service.PhotosService
}

func NewPhotoHandler(svc service.PhotosService) PhotoHandler {
	return &photoHandlerImpl{
		svc: svc,
	}
}

func (p *photoHandlerImpl) UpdatePhoto(ctx *gin.Context) {
	var data model.UpdatePhoto

	photoID, err := strconv.Atoi(ctx.Param("id"))
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
	updatedPhoto, err := p.svc.UpdatePhoto(ctx, data, photoID, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// Respond with updated photo details
	ctx.JSON(http.StatusOK, updatedPhoto)
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
func (p *photoHandlerImpl) GetAllPhotos(ctx *gin.Context) {
	photos, err := p.svc.GetAllPhotos(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, photos)
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
func (p *photoHandlerImpl) CreatePhoto(ctx *gin.Context) {
	photoCreate := model.CreatePhoto{}

	err := ctx.ShouldBindJSON(&photoCreate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid request body"})
		return
	}

	errs := photoCreate.PhotoValidate()
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

	// Call service to create photo
	photo, err := p.svc.CreatePhoto(ctx, photoCreate, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, photo)
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
func (p *photoHandlerImpl) DeletePhoto(ctx *gin.Context) {
	// Get photo ID
	photoID, err := strconv.Atoi(ctx.Param("id"))
	if photoID == 0 || err != nil {
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

	// Check if the photo belongs to the user
	if photoID != int(userId) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid request: photo does not belong to the user"})
		return
	}

	// Call service to delete photo
	err = p.svc.DeletePhoto(ctx, photoID, int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
