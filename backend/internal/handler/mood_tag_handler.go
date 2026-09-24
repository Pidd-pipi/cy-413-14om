package handler

import (
	"github.com/blueship581/mindgarden/backend/internal/constants"
	"github.com/blueship581/mindgarden/backend/internal/dto"
	"github.com/blueship581/mindgarden/backend/internal/middleware"
	"github.com/blueship581/mindgarden/backend/internal/service"
	"github.com/blueship581/mindgarden/backend/internal/util"
	"github.com/gin-gonic/gin"
	"log/slog"
	"strconv"
)

type MoodTagHandler struct {
	s      *service.MoodTagService
	logger *slog.Logger
}

func NewMoodTagHandler(s *service.MoodTagService, l *slog.Logger) *MoodTagHandler {
	return &MoodTagHandler{s, l}
}
func (h *MoodTagHandler) List(c *gin.Context) {
	v, e := h.s.ListActive(middleware.UserID(c))
	if e != nil {
		c.Error(e)
		return
	}
	h.logger.Info(constants.LogMoodTagListed)
	ok(c, v)
}
func (h *MoodTagHandler) Create(c *gin.Context) {
	var r dto.MoodTagRequest
	if !bind(c, &r) {
		return
	}
	v, e := h.s.Create(middleware.UserID(c), r.Name)
	if e != nil {
		c.Error(e)
		return
	}
	created(c, v)
}
func (h *MoodTagHandler) Deactivate(c *gin.Context) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		c.Error(util.NewAppError(constants.CodeValidation, "MoodTag[id] deactivate failed: invalid id", e))
		return
	}
	if e = h.s.Deactivate(middleware.UserID(c), uint(id)); e != nil {
		c.Error(e)
		return
	}
	ok(c, gin.H{"deactivated": id})
}
