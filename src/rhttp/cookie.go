package rhttp

var Cookie string

func GetCookie() string {
	if Cookie != "" {
		return Cookie
	}

	return ""
}

func SetCookie(NewCookie string) {
	println("new cookie set")
	Cookie = NewCookie
}
