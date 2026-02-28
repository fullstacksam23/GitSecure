package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/fullstacksam23/GitSecure/internal/api"
)

func main() {
	app := api.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel() // cancels the current context at the end of the current function

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("server failed to start")
	}
}
