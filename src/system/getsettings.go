package system

import (
	"encoding/json"
	"os"
	"path"
)

type SettingsStruct struct {
	StartonStartup bool
	ShowConsole    bool
}

var Settings SettingsStruct = SettingsStruct{StartonStartup: true, ShowConsole: false}
var ProgramDirectory string

const (
	OS_READ       = 04
	OS_WRITE      = 02
	OS_USER_SHIFT = 6

	OS_USER_R  = OS_READ << OS_USER_SHIFT
	OS_USER_W  = OS_WRITE << OS_USER_SHIFT
	OS_USER_RW = OS_USER_R | OS_USER_W
)

func GetAppdata() string {
	Appdata, err := os.UserConfigDir()

	if err != nil {
		println(err.Error())
		return ""
	}

	return Appdata
}

func GetProgramDirectory() string {
	if ProgramDirectory != "" {
		return ProgramDirectory
	}

	Appdata := GetAppdata()

	if Appdata == "" {
		return ""
	}

	if _, err := os.Stat(path.Join(Appdata, "QoLDiscordClient")); os.IsNotExist(err) {
		err := os.Mkdir(path.Join(Appdata, "QoLDiscordClient"), OS_USER_RW)
		if err != nil {
			println(err.Error())
		}
	}

	ProgramDirectory = path.Join(Appdata, "QoLDiscordClient")
	return ProgramDirectory
}

func GetSettings() {
	Directory := GetProgramDirectory()
	if Directory == "" {
		return
	}

	bytes, err := os.ReadFile(path.Join(Directory, "settings.json"))

	if err != nil {
		if os.IsNotExist(err) {
			SaveSettings()
			return
		}

		println(err.Error())
		return
	}

	json.Unmarshal(bytes, &Settings)
}

func SaveSettings() {
	Directory := GetProgramDirectory()
	if Directory == "" {
		return
	}

	bytes, err := json.Marshal(Settings)

	if err != nil {
		println(err.Error())
		return
	}

	WriteErr := os.WriteFile(path.Join(Directory, "settings.json"), bytes, OS_USER_RW)
	if WriteErr != nil {
		println(WriteErr.Error())
	}
}
