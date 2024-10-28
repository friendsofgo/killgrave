package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/friendsofgo/killgrave/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
