package router

import (
	"mygram/handler"
	"mygram/middleware"

	"github.com/gin-gonic/gin"
)

type CommentsRouter interface {
	Mount()
}

type commentsRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.CommentsHandler
}

func NewCommentsRouter(v *gin.RouterGroup, handler handler.CommentsHandler) CommentsRouter {
	return &commentsRouterImpl{v: v, handler: handler}
}

func (c *commentsRouterImpl) Mount() {
	c.v.Use(middleware.CheckAuthBearer)
	c.v.POST("", c.handler.CreateComment)
	c.v.GET("", c.handler.GetAllComment)
	c.v.DELETE("/:commentId", c.handler.DeleteComment)
	c.v.PUT("/:commentId", c.handler.UpdateComment)
}
