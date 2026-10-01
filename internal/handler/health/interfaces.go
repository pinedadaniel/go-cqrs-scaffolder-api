package health

import (
	"context"

	"github.com/pinedadaniel/go-cqrs-scaffolder-api/internal/core/usecase/health"
)

type UseCase interface {
	Execute(ctx context.Context) (health.Output, error)
}
