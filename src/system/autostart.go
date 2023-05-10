package system

import (
	"os"

	"github.com/emersion/go-autostart"
)

var exe string
var app *autostart.App

func FetchAutoStartInfo() {
	if app == nil {
		app = &autostart.App{
			Name:        "QoLDiscordClient",
			DisplayName: "QoL Discord Client for Roblox",
			Exec:        []string{exe},
			Icon:        string(TrayIcon),
		}
	}
	if exe == "" {
		exe, _ = os.Executable()
	}
}

func UpdateAutoStartState(Enabled bool) {
	FetchAutoStartInfo()
	if Enabled {
		if !app.IsEnabled() {
			app.Enable()
		}
	} else if !Enabled {
		if app.IsEnabled() {
			app.Disable()
		}
	}
}
