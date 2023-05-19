package system

import (
	"runtime"

	"github.com/gonutz/w32/v2"
)

func ShowConsole(Enabled bool) {
	if runtime.GOOS != "windows" {
		return
	}

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
