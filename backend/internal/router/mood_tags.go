package router

import (
	"github.com/blueship581/mindgarden/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterMoodTags(g *gin.RouterGroup, h *handler.MoodTagHandler, auth gin.HandlerFunc) {
	p := g.Group("/mood-tags", auth)
	p.GET("", h.List)
	p.POST("", h.Create)
	p.DELETE("/:id", h.Deactivate)
}
