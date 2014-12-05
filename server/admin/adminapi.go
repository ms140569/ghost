package admin

import (
	"github.com/ms140569/ghost/globals"
)

func ListDestinations() []string {
	return globals.Config.Storage.ListDestinations()
}

func CreateDestinattion(destination string) error {
	return nil
}

func DeleteDestination(destination string) error {
	return nil
}

func StatusDestination(destination string) (string, error) {
	return "", nil
}
