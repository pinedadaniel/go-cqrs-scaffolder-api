package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pinedadaniel/logger-go/pkg/log"
)

var (
	execCommand        = exec.Command
	chdirCommand       = os.Chdir
	local              = "local"
	ErrEmptyScope      = errors.New("SCOPE environment variable is empty")
	ErrEmptyAppCommand = errors.New("APP_COMMAND environment variable is empty")
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Panic(ctx, "could not start app",
			log.Err(err),
		)
	}
}

func run(ctx context.Context) error {
	scope, app, err := env()
	if err != nil {
		return err
	}

	if scope == local {
		log.Init(true)
	}

	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}

	log.Info(ctx, "Running bootstrapper with",
		log.String("scope", scope),
		log.String("app", app),
	)

	targetDir := filepath.Join("cmd", app)

	if errCommand := chdirCommand(targetDir); errCommand != nil {
		return fmt.Errorf("change dir error on: %s - %w", app, errCommand)
	}

	log.Info(ctx, "changing path to new dir",
		log.String("app", app),
	)
	cmd := execCommand("./app", os.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	configPath := filepath.Join(rootDir, "internal", "config", "profile")
	cmd.Env = append(os.Environ(), fmt.Sprintf("CONFIG_DIR=%s", configPath))

	log.Info(ctx, configPath)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start app: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("could not wait app: %w", err)
	}

	return nil
}

func env() (string, string, error) {
	scope := os.Getenv("SCOPE")
	if scope == "" {
		return "", "", ErrEmptyScope
	}

	scope = strings.Split(scope, "-")[0]

	app := os.Getenv("APP_COMPONENT")
	if app == "" {
		return scope, "", ErrEmptyAppCommand
	}

	app = strings.Split(app, "-")[0]

	return scope, app, nil
}
