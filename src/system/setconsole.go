package system

import "github.com/gonutz/w32/v2"

func ShowConsole(Enabled bool) {
	console := w32.GetConsoleWindow()

	if console == 0 {
		return
	}

	_, consoleProcID := w32.GetWindowThreadProcessId(console)
	if w32.GetCurrentProcessId() == consoleProcID {
		var State int

		if Enabled {
			State = w32.SW_SHOW
		} else {
			State = w32.SW_HIDE
		}

		w32.ShowWindowAsync(console, State)
	}
}
