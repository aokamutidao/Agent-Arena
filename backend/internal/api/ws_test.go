package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"agent-arena/backend/internal/engine"

	"github.com/gorilla/websocket"
)

func setupWSServer() *Server {
	server := NewServer()

	// 创建测试对局
	eng := engine.NewEngine()
	game := eng.NewGame(1, "Berserker", "Tactician", "berserker", "tactician")
	game.Status = engine.StatusPlaying
	game.CurrentRound = 5
	game.AgentRed.HP = 85
	game.AgentBlue.HP = 70
	server.store.CreateGame(game)

	return server
}

func TestWSHub_NewHub(t *testing.T) {
	hub := NewWSHub()
	if hub == nil {
		t.Fatal("hub should not be nil")
	}
	if len(hub.rooms) != 0 {
		t.Fatal("new hub should have no rooms")
	}
}

func TestWSHub_RegisterUnregister(t *testing.T) {
	hub := NewWSHub()
	go hub.Run()

	// 创建模拟客户端
	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		gameID: 1,
	}

	hub.register <- client
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount(1) != 1 {
		t.Fatalf("expected 1 client, got %d", hub.ClientCount(1))
	}

	hub.unregister <- client
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount(1) != 0 {
		t.Fatalf("expected 0 clients, got %d", hub.ClientCount(1))
	}
}

func TestWSHub_Broadcast(t *testing.T) {
	hub := NewWSHub()
	go hub.Run()

	client := &Client{
		hub:    hub,
		send:   make(chan []byte, 256),
		gameID: 1,
	}

	hub.register <- client
	time.Sleep(50 * time.Millisecond)

	// 广播消息
	data := GameStateData{
		GameID:       1,
		Status:       "playing",
		CurrentRound: 5,
		AgentRedHP:   85,
		AgentBlueHP:  70,
	}
	hub.Broadcast(1, "game_state", data)

	// 检查客户端收到消息
	select {
	case msg := <-client.send:
		var wsMsg WSMessage
		json.Unmarshal(msg, &wsMsg)
		if wsMsg.Type != "game_state" {
			t.Fatalf("expected game_state, got %s", wsMsg.Type)
		}
		var stateData GameStateData
		json.Unmarshal(wsMsg.Data, &stateData)
		if stateData.CurrentRound != 5 {
			t.Fatalf("expected round 5, got %d", stateData.CurrentRound)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for broadcast message")
	}
}

func TestWSHub_BroadcastToEmptyRoom(t *testing.T) {
	hub := NewWSHub()
	go hub.Run()

	// 广播到不存在的房间不应 panic
	data := GameStateData{GameID: 999, Status: "playing"}
	hub.Broadcast(999, "game_state", data)
	// 如果没有 panic 就通过
}

func TestWSHub_MultipleClients(t *testing.T) {
	hub := NewWSHub()
	go hub.Run()

	c1 := &Client{hub: hub, send: make(chan []byte, 256), gameID: 1}
	c2 := &Client{hub: hub, send: make(chan []byte, 256), gameID: 1}
	c3 := &Client{hub: hub, send: make(chan []byte, 256), gameID: 2}

	hub.register <- c1
	hub.register <- c2
	hub.register <- c3
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount(1) != 2 {
		t.Fatalf("expected 2 clients in room 1, got %d", hub.ClientCount(1))
	}
	if hub.ClientCount(2) != 1 {
		t.Fatalf("expected 1 client in room 2, got %d", hub.ClientCount(2))
	}

	// 广播到 room 1
	hub.Broadcast(1, "test", map[string]string{"msg": "hello"})
	time.Sleep(50 * time.Millisecond)

	// c1, c2 应该收到，c3 不应该
	if len(c1.send) != 1 {
		t.Fatalf("c1 should have 1 message, got %d", len(c1.send))
	}
	if len(c2.send) != 1 {
		t.Fatalf("c2 should have 1 message, got %d", len(c2.send))
	}
	if len(c3.send) != 0 {
		t.Fatalf("c3 should have 0 messages, got %d", len(c3.send))
	}
}

func TestWSHandler_Connect(t *testing.T) {
	server := setupWSServer()

	// 创建 httptest server
	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	// WebSocket 连接
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?game_id=1"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("ws dial error: %v", err)
	}
	defer ws.Close()

	// 等待客户端注册
	time.Sleep(200 * time.Millisecond)

	if server.wsHub.ClientCount(1) != 1 {
		t.Fatalf("expected 1 client, got %d", server.wsHub.ClientCount(1))
	}
}

func TestWSHandler_InvalidGameID(t *testing.T) {
	server := setupWSServer()

	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	// 无效 game_id
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?game_id=abc"
	resp, err := http.Get(ts.URL + "/ws?game_id=abc")
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	_ = wsURL
}

func TestWSHandler_BroadcastAfterConnect(t *testing.T) {
	server := setupWSServer()

	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws?game_id=1"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("ws dial error: %v", err)
	}
	defer ws.Close()

	// 等待客户端注册
	time.Sleep(200 * time.Millisecond)

	if server.wsHub.ClientCount(1) != 1 {
		t.Fatalf("expected 1 client, got %d", server.wsHub.ClientCount(1))
	}

	// 通过 Hub 直接广播
	turnData := TurnUpdateData{
		Round:      6,
		RedAction:  ActionBrief{Type: "ATTACK"},
		BlueAction: ActionBrief{Type: "MOVE", Target: &Pos{X: 5, Y: 5}},
		RedHP:      85,
		BlueHP:     55,
	}
	server.wsHub.Broadcast(1, "turn_update", turnData)

	// 设置读取超时并读取消息
	ws.SetReadDeadline(time.Now().Add(3 * time.Second))

	// 可能读到快照或回合更新（取决于时序），验证至少能收到一条
	_, msg, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("read message error: %v", err)
	}

	var wsMsg WSMessage
	json.Unmarshal(msg, &wsMsg)
	if wsMsg.Type != "game_state" && wsMsg.Type != "turn_update" {
		t.Fatalf("expected game_state or turn_update, got %s", wsMsg.Type)
	}
}

func TestWSTypes_TurnUpdate(t *testing.T) {
	data := TurnUpdateData{
		Round:      10,
		RedAction:  ActionBrief{Type: "MOVE", Target: &Pos{X: 3, Y: 4}},
		BlueAction: ActionBrief{Type: "SKILL"},
		RedHP:      80,
		BlueHP:     68,
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded TurnUpdateData
	json.Unmarshal(b, &decoded)

	if decoded.Round != 10 {
		t.Fatalf("expected round 10, got %d", decoded.Round)
	}
	if decoded.RedAction.Target == nil {
		t.Fatal("expected red action target")
	}
	if decoded.RedAction.Target.X != 3 {
		t.Fatalf("expected x=3, got %d", decoded.RedAction.Target.X)
	}
	if decoded.BlueAction.Target != nil {
		t.Fatal("blue action should have no target")
	}
}
