package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kiritoxkiriko/comical-tool/internal/service"
)

type ShortHandler struct {
	*Handler
	shortService service.ShortService
}

func NewShortHandler(
    handler *Handler,
    shortService service.ShortService,
) *ShortHandler {
	return &ShortHandler{
		Handler:      handler,
		shortService: shortService,
	}
}

func (h *ShortHandler) GetShort(ctx *gin.Context) {

}
