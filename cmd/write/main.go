package main

import (
	"context"

	"github.com/pinedadaniel/logger-go/pkg/log"
)

func main() {
	ctx := context.Background()

	log.Info(ctx, "PONG")
}
