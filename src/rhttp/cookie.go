package rhttp

import "github.com/Haydz6/Roblox-QoL-Discord-Client/src/system"

var Cookie string

func GetCookie() string {
	if Cookie != "" {
		return Cookie
	}

	return ""
}

func SetCookie(NewCookie string) {
	Cookie = NewCookie
	system.SaveSettings()
}
