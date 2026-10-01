package health

import (
	"context"
)

type UseCase struct {
	healthRepository Repository
	Version          string
}

func New(repo Repository, version string) *UseCase {
	return &UseCase{
		healthRepository: repo,
		Version:          version,
	}
}

func (uc *UseCase) Execute(ctx context.Context) (Output, error) {
	health, err := uc.healthRepository.GetHealth(ctx, uc.Version)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		Version:   health.Version,
		Timestamp: health.Timestamp,
		Status:    string(health.Status),
	}

	return output, nil
}
