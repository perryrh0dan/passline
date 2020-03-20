package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()

	//trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	app := setupApp(ctx)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
