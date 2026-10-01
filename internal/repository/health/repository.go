package health

import (
	"context"
	"time"

	"github.com/pinedadaniel/scaffolder-api-go/internal/core/domain"
)

type Repository struct {
}

func New() (*Repository, error) {
	return &Repository{}, nil
}

func (m *Repository) GetHealth(ctx context.Context, version string) (domain.Health, error) {
	return domain.Health{
		Status:    domain.StatusUP,
		Timestamp: time.Now().UTC(),
		Version:   version,
	}, nil
}
