package admin

import (
	"github.com/ms140569/ghost/globals"
)

func ListDestinations() []string {
	return globals.Config.Storage.ListDestinations()
}

func CreateDestinattion(destination string) {}

func DeleteDestination(destination string) {}

func StatusDestination(destination string) string {
	return ""
}
