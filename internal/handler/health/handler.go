package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/logger-go/pkg/log"
	"github.com/pinedadaniel/scaffolder-api-go/internal/handler"
)

type Handler struct {
	healthUseCase UseCase
}

func New(useCase UseCase) *Handler {
	return &Handler{
		healthUseCase: useCase,
	}
}

func (h *Handler) Get(ctx *gin.Context) {

	log.Info(ctx, "Request Init")
	health, err := h.healthUseCase.Execute(ctx.Request.Context())
	if err != nil {
		log.Error(ctx, "Error HealthUseCase", log.Err(err))
		handleError(ctx, err)
	}

	res := Response{
		Status:    health.Status,
		Timestamp: health.Timestamp,
		Version:   health.Version,
	}

	ctx.JSON(http.StatusOK, handler.Success(res))
}
