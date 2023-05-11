package rhttp

var Cookie string

func GetCookie() string {
	if Cookie != "" {
		return Cookie
	}

	return ""
}

func SetCookie(NewCookie string) {
	Cookie = NewCookie
}
