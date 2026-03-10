package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/fullstacksam23/GitSecure/internal/api"
	"github.com/fullstacksam23/GitSecure/internal/worker"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel() // cancels the current context at the end of the current function

	mode := os.Getenv("MODE")

	if mode == "worker" {
		fmt.Println("Starting worker...")
		worker.StartWorker(ctx)
		return
	}

	fmt.Println("Starting API server")
	app := api.New()
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("server failed to start")
	}
}
