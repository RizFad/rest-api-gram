package router

import (
	"mygram/handler"
	"mygram/middleware"

	"github.com/gin-gonic/gin"
)

type PhotoRouter interface {
	Mount()
}

type photoRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.PhotoHandler
}

func NewPhotoRouter(v *gin.RouterGroup, handler handler.PhotoHandler) PhotoRouter {
	return &photoRouterImpl{v: v, handler: handler}
}

func (p *photoRouterImpl) Mount() {
	p.v.Use(middleware.CheckAuthBearer)
	p.v.POST("", p.handler.CreatePhoto)
	p.v.GET("", p.handler.GetAllPhotos)
	p.v.DELETE("/:id", p.handler.DeletePhoto)
	p.v.PUT("/:id", p.handler.UpdatePhoto)
}
