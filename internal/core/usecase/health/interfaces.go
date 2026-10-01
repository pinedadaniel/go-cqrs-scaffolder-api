package health

import (
	"context"

	"github.com/pinedadaniel/go-cqrs-scaffolder-api/internal/core/domain"
)

type Repository interface {
	GetHealth(ctx context.Context, version string) (domain.Health, error)
}
