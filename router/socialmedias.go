package router

import (
	"mygram/handler"
	"mygram/middleware"

	"github.com/gin-gonic/gin"
)

type SocialMediasRouter interface {
	Mount()
}

type socialmediasRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.SocialMediasHandler
}

func NewSocialMediasRouter(v *gin.RouterGroup, handler handler.SocialMediasHandler) SocialMediasRouter {
	return &socialmediasRouterImpl{v: v, handler: handler}
}

func (sm *socialmediasRouterImpl) Mount() {
	sm.v.Use(middleware.CheckAuthBearer)
	sm.v.POST("", sm.handler.CreateSocialMedia)
	sm.v.GET("", sm.handler.GetAllSocialMedia)
	sm.v.DELETE("/:socialMediaId", sm.handler.DeleteSocialMedia)
	sm.v.PUT("/:socialMediaId", sm.handler.UpdateSocialMedia)
}
