package system

import (
	"os"

	_ "embed"

	"github.com/Haydz6/rich-go/client"
	"github.com/getlantern/systray"
)

//go:embed RobloxFavicon.ico
var TrayIcon []byte

func CreateTray() {
	systray.Run(TrayReady, TrayFail)
}

func TrayReady() {
	systray.SetIcon(TrayIcon)
	systray.SetTitle("QoL Discord Client")
	systray.SetTooltip("QoL Discord Client")

	QuitButton := systray.AddMenuItem("Quit", "Turns off the discord client")
	StartupButton := systray.AddMenuItemCheckbox("Launch on startup", "If the program should start when you first login", Settings.StartonStartup)

	for {
		select {
		case <-QuitButton.ClickedCh:
			quit()
		case <-StartupButton.ClickedCh:
			Settings.StartonStartup = !Settings.StartonStartup
			if Settings.StartonStartup {
				StartupButton.Check()
			} else {
				StartupButton.Uncheck()
			}
			SaveSettings()
		}
	}
}

func TrayFail() {}

func quit() {
	client.Logout()
	systray.Quit()
	os.Exit(0)
}
