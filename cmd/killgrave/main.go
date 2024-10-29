package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/friendsofgo/killgrave/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
