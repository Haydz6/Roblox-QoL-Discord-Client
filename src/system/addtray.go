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

	StartupButton := systray.AddMenuItemCheckbox("Launch on startup", "If the program should start when you first login", Settings.StartonStartup)
	//ConsoleButton := systray.AddMenuItemCheckbox("Show console", "Shows debugging console", Settings.ShowConsole)
	QuitButton := systray.AddMenuItem("Quit", "Turns off the discord client")

	for {
		select {
		case <-QuitButton.ClickedCh:
			quit()
		case <-StartupButton.ClickedCh:
			Settings.StartonStartup = !Settings.StartonStartup
			UpdateAutoStartState(Settings.StartonStartup)
			if Settings.StartonStartup {
				StartupButton.Check()
			} else {
				StartupButton.Uncheck()
			}
			SaveSettings()
			// case <-ConsoleButton.ClickedCh:
			// 	Settings.ShowConsole = !Settings.ShowConsole
			// 	ShowConsole(Settings.ShowConsole)
			// 	if Settings.ShowConsole {
			// 		ConsoleButton.Check()
			// 	} else {
			// 		ConsoleButton.Uncheck()
			// 	}
			// 	SaveSettings()
		}
	}
}

func TrayFail() {}

func quit() {
	client.Logout()
	systray.Quit()
	os.Exit(0)
}
