package presence

import (
	"encoding/json"

	"github.com/Haydz6/Roblox-QoL-Discord-Client/src/rhttp"
)

var UserId int

type AuthenticatedStruct struct {
	Id int `json:"id"`
}

func GetUserId() (bool, int) {
	if UserId != 0 {
		return true, UserId
	}

	Success, Result := rhttp.RobloxRequest("https://users.roblox.com/v1/users/authenticated", "GET", nil, "", true)
	if Result != nil {
		defer Result.Body.Close()
	}

	if !Success {
		return false, 0
	}

	var Body AuthenticatedStruct
	err := json.NewDecoder(Result.Body).Decode(&Body)

	if err != nil {
		return false, 0
	}

	UserId = Body.Id
	return true, UserId
}
