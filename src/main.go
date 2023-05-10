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
	Activity *client.Activity
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

	time.Sleep(time.Second * 5)

	if len(s.conns) == 0 {
		client.SetActivity(client.Activity{State: "end"})
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	Cleared := false
	go func() {
		for range time.Tick(time.Second * 1) {
			_, err := ws.Write([]byte("ping"))
			if err != nil {
				Cleared = true
				s.cleanWS(ws)
			}
		}
	}()

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

		var Error error
		if Body.Activity != nil {
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
	go system.UpdateAutoStartState(true)
	go system.CreateTray()
	server := NewServer()

	http.HandleFunc("/presence", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(server.handleWS)}
		s.ServeHTTP(w, req)
	})

	client.Login("1105722413905346660")

	for i := 0; i <= 4; i++ {
		err := http.ListenAndServe(":"+strconv.Itoa(9300+i), nil)
		if err == nil {
			break
		}
	}
}
