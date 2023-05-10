package presence

import (
	"fmt"
	"time"

	"github.com/Haydz6/Roblox-QoL-Discord-Client/src/rhttp"
	"github.com/Haydz6/rich-go/client"
)

var DependentPresenceEnabled = false
var DependentPresenceTick = 0

func RunExternalPresence() {
	var LastPlaceId int
	var CachedDependentPresenceTick = DependentPresenceTick

	for range time.Tick(time.Second * 10) {
		if CachedDependentPresenceTick != DependentPresenceTick {
			println("tick difference")
			break
		}

		if rhttp.GetCookie() == "" {
			client.SetActivity(client.Activity{State: "end"})
			SetDependentPresence(false)
			break
		}

		Success, UserId := GetUserId()

		if !Success {
			continue
		}

		Success, InGame, UniverseId, PlaceId, JobId := GetPresence(UserId)

		if !Success {
			continue
		}

		if PlaceId == LastPlaceId {
			continue
		}

		if InGame {
			Success, PlaceInfo := GetPlaceInfo(UniverseId)
			if !Success {
				continue
			}

			ThumbnailURL := GetPlaceThumbnail(UniverseId)
			var Verified string

			if PlaceInfo.Creator.HasVerifiedBadge {
				Verified = " ☑️"
			}

			LastPlaceId = PlaceId

			GameElapsed := time.Now()
			client.SetActivity(client.Activity{
				Details:    PlaceInfo.Name,
				Buttons:    []*client.Button{{Label: "Join", Url: fmt.Sprintf("roblox://experiences/start?placeId=%d&gameInstanceId=%s", PlaceId, JobId)}, {Label: "View Game", Url: fmt.Sprintf("https://www.roblox.com/games/%d", PlaceId)}},
				State:      fmt.Sprintf("by %s%s", PlaceInfo.Creator.Name, Verified),
				LargeText:  PlaceInfo.Name,
				LargeImage: ThumbnailURL,
				SmallText:  "Roblox",
				SmallImage: "https://cdn.discordapp.com/app-assets/1105722413905346660/1105722508038115438.png",
				Timestamps: &client.Timestamps{Start: &GameElapsed},
			})
		} else {
			LastPlaceId = 0
			client.SetActivity(client.Activity{State: "end"})
		}
	}
}

func SetDependentPresence(Enabled bool) bool {
	if DependentPresenceEnabled == Enabled {
		return DependentPresenceEnabled
	}

	HasCookie := rhttp.GetCookie() != ""
	DependentPresenceTick++
	if Enabled && HasCookie {
		println("in external mode")
		go RunExternalPresence()
	}

	return HasCookie
}
