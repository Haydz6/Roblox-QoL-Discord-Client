package rhttp

import (
	"bytes"
	"fmt"
	"net/http"
)

type ResponseStruct struct {
	Success bool
	Result  string
}

var CSRFToken string

func RobloxRequest(URL string, Method string, Headers map[string]string, Body string, RequiresAuth bool) (bool, *http.Response) {
	if GetCookie() == "" {
		return false, nil
	}

	client := &http.Client{}
	request, _ := http.NewRequest(Method, URL, bytes.NewBuffer([]byte(Body)))

	if Headers != nil {
		for k, v := range Headers {
			request.Header.Set(k, v)
		}
	}

	if Headers == nil || Headers["Content-Type"] == "" {
		request.Header.Set("Content-Type", "application/json")
	}
	request.Header.Set("x-csrf-token", CSRFToken)

	if RequiresAuth {
		request.Header.Set("Cookie", fmt.Sprintf(".ROBLOSECURITY=%s;", GetCookie()))
	}

	response, err := client.Do(request)

	if err == nil {
		if response.StatusCode == 403 {
			CSRFToken = response.Header.Get("x-csrf-token")

			if CSRFToken != "" {
				return RobloxRequest(URL, Method, Headers, Body, RequiresAuth)
			}
		} else if response.StatusCode == 401 {
			SetCookie("") //Cookie is invalid
			return false, nil
		}
	} else {
		println(err.Error())
	}

	return err == nil && response.StatusCode < 400, response
}
