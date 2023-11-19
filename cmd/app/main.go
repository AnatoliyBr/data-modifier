package main

import (
	"log"

	"github.com/AnatoliyBr/data-modifier/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
