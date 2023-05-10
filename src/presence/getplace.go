package presence

import (
	"encoding/json"
	"fmt"

	"github.com/Haydz6/Roblox-QoL-Discord-Client/src/rhttp"
)

type GetUniverseStruct struct {
	Id          int    `json:"id"`
	RootPlaceId int    `json:"rootPlaceId"`
	Name        string `json:"name"`
	Creator     struct {
		Name             string `json:"name"`
		HasVerifiedBadge bool   `json:"hasVerifiedBadge"`
	} `json:"creator"`
}

type UniverseResultStruct struct {
	Data []GetUniverseStruct `json:"data"`
}

type ThumbnailResultStruct struct {
	Data []struct {
		ImageUrl string `json:"imageUrl"`
	} `json:"data"`
}

var CachedGetPlace *GetUniverseStruct
var CachedThumbnailURL string
var CachedThumbnailId int

func GetPlaceInfo(UniverseId int) (bool, *GetUniverseStruct) {
	if CachedGetPlace != nil && CachedGetPlace.Id == UniverseId {
		return true, CachedGetPlace
	}

	Success, Result := rhttp.RobloxRequest(fmt.Sprintf("https://games.roblox.com/v1/games?universeIds=%d", UniverseId), "GET", nil, "", true)
	if Result != nil {
		defer Result.Body.Close()
	}

	if !Success {
		return false, nil
	}

	var Body UniverseResultStruct
	JSONErr := json.NewDecoder(Result.Body).Decode(&Body)

	if JSONErr != nil {
		return false, nil
	}

	CachedGetPlace = &Body.Data[0]
	return true, CachedGetPlace
}

func GetPlaceThumbnail(UniverseId int) string {
	if CachedThumbnailId == UniverseId {
		return CachedThumbnailURL
	}

	Success, Result := rhttp.RobloxRequest(fmt.Sprintf("https://thumbnails.roblox.com/v1/games/icons?universeIds=%d&returnPolicy=PlaceHolder&size=512x512&format=Png&isCircular=false", UniverseId), "GET", nil, "", true)
	if Result != nil {
		defer Result.Body.Close()
	}

	if !Success {
		return "https://tr.rbxcdn.com/53eb9b17fe1432a809c73a13889b5006/512/512/Image/Png"
	}

	var Body ThumbnailResultStruct
	JSONErr := json.NewDecoder(Result.Body).Decode(&Body)

	if JSONErr != nil {
		return "https://tr.rbxcdn.com/53eb9b17fe1432a809c73a13889b5006/512/512/Image/Png"
	}

	CachedThumbnailId = UniverseId
	CachedThumbnailURL = Body.Data[0].ImageUrl

	return CachedThumbnailURL
}
