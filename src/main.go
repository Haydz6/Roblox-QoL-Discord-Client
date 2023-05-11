package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Haydz6/Roblox-QoL-Discord-Client/src/presence"
	"github.com/Haydz6/Roblox-QoL-Discord-Client/src/rhttp"
	"github.com/Haydz6/Roblox-QoL-Discord-Client/src/system"
	"github.com/Haydz6/rich-go/client"
	"golang.org/x/net/websocket"
)

type PresenceUpdate struct {
	Activity       *client.Activity
	PlaceId        int
	JobId          string
	Authentication string
}

type MessageToBrowserStruct struct {
	Type       string `json:",omitempty"`
	Timestamp  int64  `json:",omitempty"`
	PlaceId    int    `json:",omitempty"`
	UniverseId int    `json:",omitempty"`
	JobId      string `json:",omitempty"`
}

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	for Oldws := range s.conns {
		Oldws.Close()
	}

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) cleanWS(ws *websocket.Conn) {
	delete(s.conns, ws)

	println("Setting?")
	if !presence.SetDependentPresence(true) {
		println("No cookie!")
		time.Sleep(time.Second * 5)

		println(len(s.conns))
		if len(s.conns) == 0 {
			client.SetActivity(client.Activity{State: "end"})
		}
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	Cleared := false
	go func() {
		for range time.Tick(time.Second * 1) {
			_, err := ws.Write([]byte(""))
			if err != nil {
				Cleared = true
				s.cleanWS(ws)
				break
			}
		}
	}()

	presence.SetDependentPresence(false)
	println(presence.LastPlaceId)
	if presence.LastPlaceId != 0 {
		bytes, err := json.Marshal(MessageToBrowserStruct{Type: "Timestamp", Timestamp: presence.LastTimestamp.UnixMilli(), PlaceId: presence.LastPlaceId, JobId: presence.LastJobId, UniverseId: presence.LastUniverseId})
		if err == nil {
			ws.Write(bytes)
		}
	}

	for !Cleared {
		n, err := ws.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		msg := buf[:n]
		var Body PresenceUpdate
		JSONErr := json.Unmarshal(msg, &Body)

		if JSONErr != nil {
			ws.Write([]byte(JSONErr.Error()))
			continue
		}

		if Body.Authentication != "" {
			rhttp.SetCookie(Body.Authentication)
			continue
		}

		var Error error
		if Body.Activity != nil {
			presence.LastJobId = Body.JobId
			presence.LastPlaceId = Body.PlaceId
			presence.LastTimestamp = *Body.Activity.Timestamps.Start
			Error = client.SetActivity(*Body.Activity)
		} else {
			Error = client.SetActivity(client.Activity{State: "end"})
		}

		if Error != nil {
			ws.Write([]byte(Error.Error()))
		}
	}
}

func main() {
	system.GetSettings()
	if system.Settings.StartonStartup {
		go system.UpdateAutoStartState(true)
	}
	go system.CreateTray()
	presence.SetDependentPresence(true)
	server := NewServer()

	println("opening server")

	http.HandleFunc("/presence", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(server.handleWS)}
		s.ServeHTTP(w, req)
	})

	println("client login")
	client.Login("1105722413905346660")
	println("client logged in")

	for i := 0; i <= 4; i++ {
		err := http.ListenAndServe(":"+strconv.Itoa(9300+i), nil)
		if err == nil {
			break
		}
	}
}
