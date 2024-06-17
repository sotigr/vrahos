package vrahos

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const (
	NO_LOBBY_KEY = "none"
)

var newLineRe = regexp.MustCompile(`\r?\n`)

func generateUuid() (string, error) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}

type Lobby struct {
	mu          sync.Mutex
	Name        string
	Key         *string
	Connections []ConnectionSession
}

type ConnectionSession struct {
	id string
	w  http.ResponseWriter
	r  *http.Request
}

type EventEmitter interface {
	Init(path string)
	Send(lobby string, event_name string, data string) bool
	CreateLobby(name string, key *string)
	Close()
}

type Sse struct {
	mu      sync.Mutex
	path    string
	Lobbies *[]*Lobby
	stop    bool
}

func NewSee(path string) *Sse {
	return &Sse{
		path:    path,
		Lobbies: &[]*Lobby{},
		stop:    false,
	}
}

func (e *Sse) Init(server *http.ServeMux) {

	server.HandleFunc(e.path, func(w http.ResponseWriter, r *http.Request) {

		id, err := generateUuid()
		if err != nil {
			panic("failed to generate uuid")
		}
		lobbiesStr := r.URL.Query().Get("lobbies")
		if lobbiesStr == "" {
			io.WriteString(w, "no_lobbies")
			return
		}
		if e.stop {
			io.WriteString(w, "closed")
			return
		}

		lobbies := strings.Split(lobbiesStr, ",")
		if len(lobbies) == 0 {

			io.WriteString(w, "no_lobbies")
			return
		}
		subbed_lobbies := []*Lobby{}

		for _, lobbyStr := range lobbies {
			values := strings.Split(lobbyStr, ":")
			lobbyName := values[0]
			var lobbyKey string
			if len(values) > 1 {
				lobbyKey = values[1]
			}
			// fmt.Println(this.Lobbies)
			cur_lobby := e.get_lobby_by_name(lobbyName)
			subbed_lobbies = append(subbed_lobbies, cur_lobby)
			if cur_lobby == nil {
				io.WriteString(w, "lobby_not_found, lobby "+lobbyName)
				return
			}

			if cur_lobby.Key != nil {
				if *cur_lobby.Key != NO_LOBBY_KEY && *cur_lobby.Key != lobbyKey {
					io.WriteString(w, "forbidden, lobby "+lobbyName)
					return
				}
			}
			cur_lobby.mu.Lock()
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			} else {
				log.Println("Damn, no flush")
			}

			cur_lobby.Connections = append(cur_lobby.Connections, ConnectionSession{w: w, r: r, id: id})
			cur_lobby.mu.Unlock()
		}

		// fmt.Println("new connection")
		ctx := r.Context()

		select {
		case <-ctx.Done():
			// If the request gets cancelled, log it
			// to STDERR
			e._send_to_connection(w, "close", "ok", false)

			for _, l := range subbed_lobbies {

				l.mu.Lock()
				clean_connections := []ConnectionSession{}
				for _, c := range l.Connections {
					if c.id != id {
						clean_connections = append(clean_connections, c)
					}
				}
				l.Connections = clean_connections
				l.mu.Unlock()
			}

			// fmt.Println("connection closed")
		}

	})
}

func (e *Sse) GetPath() string {
	return e.path
}

func (e *Sse) get_lobby_by_name(name string) *Lobby {
	for _, lobby := range *e.Lobbies {
		if lobby.Name == name {
			return lobby
		}
	}
	return nil
}

func (e *Sse) _send_to_connection(w http.ResponseWriter, event_name string, data string, flush bool) {
	e.mu.Lock()
	id, err := generateUuid()
	if err != nil {
		panic("failed to generate uuid")
	}

	w.Write([]byte(fmt.Sprintf("id: %s\n", id)))
	w.Write([]byte("event: " + event_name + "\n"))
	w.Write([]byte(fmt.Sprintf("data: %s\n", data)))
	w.Write([]byte("\n"))
	if flush {
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		} else {
			log.Println("Damn, no flush")
		}
	}
	e.mu.Unlock()
}

func (e *Sse) CreateLobby(name string, key string) {

	*e.Lobbies = append(*e.Lobbies, &Lobby{Name: name, Key: &key, Connections: []ConnectionSession{}})

}

func (e *Sse) Send(lobbyName string, event_name string, data string) {
	lobby := e.get_lobby_by_name(lobbyName)
	// fmt.Println("sending")
	data = newLineRe.ReplaceAllString(data, " ")
	if lobby != nil {
		for _, c := range lobby.Connections {
			e._send_to_connection(c.w, lobbyName+"."+event_name, data, true)
		}
	}
}

func CloseSse(e *Sse) {
	e.stop = true
}
