package poker

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func mustMakePlayerServer(t *testing.T) *PlayServer {
	server, err := NewPlayerServer()
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func TestGame(t *testing.T) {
	t.Run("", func(t *testing.T) {
		winner := "Michael"
		server := httptest.NewServer(mustMakePlayerServer(t))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("could not open a ws connect on %s %v ", wsURL, err)
		}
		defer ws.Close()

		writeMessage(t, ws, winner)
		fmt.Println(wsURL)
		time.Sleep(10 * time.Millisecond)
		AssertPlayerWin(t, winner)
	})
}

func writeMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Errorf("could not send message over ws connection %v", err)
	}
}
