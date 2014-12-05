package admin

import (
	"github.com/ms140569/ghost/globals"
)

func ListDestinations() []string {
	return globals.Config.Storage.ListDestinations()
}

func CreateDestination(destination string) error {
	return globals.Config.Storage.CreateDestination(destination)
}

func DeleteDestination(destination string) error {
	return globals.Config.Storage.DeleteDestination(destination)
}

func StatusDestination(destination string) (string, error) {
	return globals.Config.Storage.StatusDestination(destination)
}
