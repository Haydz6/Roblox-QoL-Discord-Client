package presence

import (
	"encoding/json"

	"github.com/Haydz6/Roblox-QoL-Discord-Client/src/rhttp"
)

type PresenceStruct struct {
	UserPresences []struct {
		UserPresenceType int    `json:"userPresenceType"`
		RootPlaceId      int    `json:"rootPlaceId"`
		UniverseId       int    `json:"universeId"`
		GameId           string `jsom:"gameId"`
	} `json:"userPresences"`
}

type PresenceSendStruct struct {
	UserIds []int `json:"userIds"`
}

func GetPresence(UserId int) (bool, bool, int, int, string) {
	RequestBytes, err := json.Marshal(PresenceSendStruct{UserIds: []int{UserId}})

	if err != nil {
		return false, false, 0, 0, ""
	}

	println(string(RequestBytes))
	Success, Result := rhttp.RobloxRequest("https://presence.roblox.com/v1/presence/users", "POST", nil, string(RequestBytes), true)
	if Result != nil {
		defer Result.Body.Close()
	}

	println(Success)
	if !Success {
		return false, false, 0, 0, ""
	}

	var Body PresenceStruct
	JSONErr := json.NewDecoder(Result.Body).Decode(&Body)

	if JSONErr != nil {
		println(JSONErr.Error())
		return false, false, 0, 0, ""
	}

	Presence := Body.UserPresences[0]
	return true, Presence.UserPresenceType == 2, Presence.UniverseId, Presence.RootPlaceId, Presence.GameId
}
