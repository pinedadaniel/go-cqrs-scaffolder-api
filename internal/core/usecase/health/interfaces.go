package health

import (
	"context"

	"github.com/pinedadaniel/scaffolder-api-go/internal/core/domain"
)

type Repository interface {
	GetHealth(ctx context.Context, version string) (domain.Health, error)
}
