package health

import (
	"context"

	"github.com/pinedadaniel/scaffolder-api-go/internal/core/usecase/health"
)

type UseCase interface {
	Execute(ctx context.Context) (health.Output, error)
}
