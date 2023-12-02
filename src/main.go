package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

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
	Type       string                      `json:",omitempty"`
	Timestamp  int64                       `json:",omitempty"`
	PlaceId    int                         `json:",omitempty"`
	UniverseId int                         `json:",omitempty"`
	JobId      string                      `json:",omitempty"`
	User       *client.AuthenticatedStruct `json:",omitempty"`
}

var LastTimestamp time.Time
var LastPlaceId int
var LastUniverseId int
var LastJobId string

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

	time.Sleep(time.Second * 5)

	if len(s.conns) == 0 {
		client.SetActivity(nil)
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

	if LastPlaceId != 0 {
		bytes, err := json.Marshal(MessageToBrowserStruct{Type: "Timestamp", Timestamp: LastTimestamp.UnixMilli(), PlaceId: LastPlaceId, JobId: LastJobId, UniverseId: LastUniverseId})
		if err == nil {
			ws.Write(bytes)
		}
	}

	println("is authed", client.Authentication != nil)
	if client.Authentication != nil {
		bytes, err := json.Marshal(MessageToBrowserStruct{Type: "Authentication", User: client.Authentication})
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
			continue
		}

		var Error error
		if Body.Activity != nil {
			LastJobId = Body.JobId
			LastPlaceId = Body.PlaceId
			LastTimestamp = *Body.Activity.Timestamps.Start
			Error = client.SetActivity(Body.Activity)
		} else {
			Error = client.SetActivity(nil)
		}

		if Error != nil {
			ws.Write([]byte(Error.Error()))
		}
	}
}

func (s *Server) SendNewAuthentication() {
	for {
		client.AuthenticationUpdate.Add(1)
		client.AuthenticationUpdate.Wait()
		println("auth update")

		bytes, err := json.Marshal(MessageToBrowserStruct{Type: "Authentication", User: client.Authentication})
		if err == nil {
			for ws := range s.conns {
				ws.Write(bytes)
			}
		}
	}
}

func main() {
	system.GetSettings()
	system.SaveSettings() //update struct
	system.ShowConsole(system.Settings.ShowConsole)
	if system.Settings.StartonStartup {
		go system.UpdateAutoStartState(true)
	}
	go system.CreateTray()

	server := NewServer()

	println("opening server")

	http.HandleFunc("/presence", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(server.handleWS)}
		s.ServeHTTP(w, req)
	})

	go server.SendNewAuthentication()

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
