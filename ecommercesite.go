package main

import (
	"context"
	"ecommercesite/application"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	appStartupErr := app.Start(ctx)

	if appStartupErr != nil {
		fmt.Println("App failed to start: %w", appStartupErr)
	}
}
