package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pinedadaniel/logger-go/pkg/log"
	"github.com/pinedadaniel/scaffolder-api-go/internal/config"
	"github.com/pinedadaniel/scaffolder-api-go/internal/core/usecase/health"
	healthHandler "github.com/pinedadaniel/scaffolder-api-go/internal/handler/health"
	healthRepository "github.com/pinedadaniel/scaffolder-api-go/internal/repository/health"
	"github.com/pinedadaniel/scaffolder-api-go/internal/routes/read"
)

var (
	local         = "local"
	ErrEmptyScope = errors.New("SCOPE environment variable is empty")
)

func main() {
	ctx := context.Background()

	scope, err := getScope()
	if err != nil {
		log.Panic(ctx, "Scope not found",
			log.Err(err),
		)
	}

	if scope == local {
		log.Init(true)
	}
	log.Info(ctx, "starting app",
		log.String("appComponent", "read"),
	)

	if err := run(ctx); err != nil {
		log.Panic(ctx, "could not start app",
			log.String("appComponent", "read"),
			log.Err(err),
		)
	}
}

func run(ctx context.Context) error {
	app, err := initializeApp()

	cfg, err := config.New(config.FileReaderImp)
	if err != nil {
		log.Error(ctx, "error reading config")
		return fmt.Errorf("could not get config: %w", err)
	}

	deps, err := buildDependencies(cfg)
	if err != nil {
		return err
	}

	read.Register(app, read.Handlers{
		Health: healthHandler.New(deps.useCases.health),
	})

	log.Info(ctx, "Server listening",
		log.String("Port", cfg.HTTP.Port))
	return app.Run(":" + cfg.HTTP.Port)
}

func initializeApp() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Recovery())

	return app, nil
}

type repositories struct {
	health *healthRepository.Repository
}

func initializeRepositories() (repositories, error) {
	healthRepo, err := healthRepository.New()

	if err != nil {
		return repositories{}, fmt.Errorf("could not initialize health repository: %w", err)
	}

	return repositories{
		health: healthRepo,
	}, nil
}

type useCases struct {
	health *health.UseCase
}

func initializeUseCases(repos repositories, cfg config.Configuration) useCases {
	return useCases{
		health: health.New(repos.health, cfg.App.Version),
	}
}

type dependencies struct {
	repositories repositories
	useCases     useCases
}

func buildDependencies(cfg config.Configuration) (dependencies, error) {
	repos, err := initializeRepositories()
	if err != nil {
		return dependencies{}, fmt.Errorf("could not initialize repositories: %w", err)
	}

	uscs := initializeUseCases(repos, cfg)

	return dependencies{
		repositories: repos,
		useCases:     uscs,
	}, nil
}

func getScope() (string, error) {
	scope := os.Getenv("SCOPE")
	if scope == "" {
		return "", ErrEmptyScope
	}
	return scope, nil
}
