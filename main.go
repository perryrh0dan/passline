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
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
			os.Exit(1)
		case <-ctx.Done():
			os.Exit(1)
		}
	}()

	app := setupApp(ctx)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
