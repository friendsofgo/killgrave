package main

import (
	"github.com/friendsofgo/killgrave/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
