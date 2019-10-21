package main

import (
	"log"
	"os"
)

func main() {
	app := setupApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
