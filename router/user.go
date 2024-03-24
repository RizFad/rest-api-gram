package router

import (
	"mygram/handler"
	"mygram/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter interface {
	Mount()
}

type userRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.UserHandler
}

func NewUserRouter(v *gin.RouterGroup, handler handler.UserHandler) UserRouter {
	return &userRouterImpl{v: v, handler: handler}
}

func (u *userRouterImpl) Mount() {
	// activity
	// /users/sign-up
	u.v.POST("/sign-up", u.handler.UserSignUp)
	u.v.POST("/login", u.handler.UserSignIn)

	// users
	u.v.Use(middleware.CheckAuthBearer)
	// /users
	u.v.GET("", u.handler.GetUsers)
	// /users/:id
	u.v.GET("/:userId", u.handler.GetUsersById)
	u.v.DELETE("/:userId", u.handler.DeleteUsersById)
	u.v.PUT("/:userId", u.handler.UpdateUsersById)
}
