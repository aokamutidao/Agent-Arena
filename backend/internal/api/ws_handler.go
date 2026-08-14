package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // MVP: 允许所有来源
	},
}

// WSHandler WebSocket 连接处理器
func (s *Server) WSHandler(c *gin.Context) {
	gameIDStr := c.Query("game_id")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil || gameID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid game_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:    s.wsHub,
		conn:   conn,
		send:   make(chan []byte, 256),
		gameID: gameID,
	}

	s.wsHub.register <- client

	// 发送当前游戏状态快照
	s.sendGameSnapshot(client, gameID)

	// 启动读写循环
	go client.writePump()
	go client.readPump()
}

// sendGameSnapshot 发送游戏状态快照给新连接的客户端
func (s *Server) sendGameSnapshot(client *Client, gameID uint64) {
	record, err := s.store.GetGame(gameID)
	if err != nil {
		return
	}

	game := record.Game
	snapshot := GameStateData{
		GameID:        game.GameID,
		Status:        string(game.Status),
		CurrentRound:  game.CurrentRound,
		AgentRedHP:    game.AgentRed.HP,
		AgentBlueHP:   game.AgentBlue.HP,
		AgentRedPosX:  game.AgentRed.Position.X,
		AgentRedPosY:  game.AgentRed.Position.Y,
		AgentBluePosX: game.AgentBlue.Position.X,
		AgentBluePosY: game.AgentBlue.Position.Y,
	}

	s.wsHub.Broadcast(gameID, "game_state", snapshot)
}
