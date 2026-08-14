package api

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client WebSocket 客户端
type Client struct {
	hub    *WSHub
	conn   *websocket.Conn
	send   chan []byte
	gameID uint64
}

// WSHub WebSocket 管理中心
type WSHub struct {
	rooms      map[uint64]map[*Client]bool // gameID -> clients
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewWSHub 创建 WebSocket Hub
func NewWSHub() *WSHub {
	return &WSHub{
		rooms:      make(map[uint64]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 启动 Hub 主循环
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.gameID] == nil {
				h.rooms[client.gameID] = make(map[*Client]bool)
			}
			h.rooms[client.gameID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.gameID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.rooms, client.gameID)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast 广播消息到指定房间
func (h *WSHub) Broadcast(gameID uint64, msgType string, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ws broadcast marshal error: %v", err)
		return
	}

	msg := WSMessage{
		Type: msgType,
		Data: dataBytes,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ws broadcast marshal error: %v", err)
		return
	}

	h.mu.RLock()
	clients, ok := h.rooms[gameID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	// 复制一份避免持锁发送
	targets := make([]*Client, 0, len(clients))
	for c := range clients {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	for _, client := range targets {
		select {
		case client.send <- msgBytes:
		default:
			// 客户端慢，跳过
			go func(c *Client) {
				h.unregister <- c
			}(client)
		}
	}
}

// BroadcastAll 广播消息到所有已连接的客户端（用于全局事件如余额更新）
func (h *WSHub) BroadcastAll(msgType string, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("ws broadcastAll marshal error: %v", err)
		return
	}

	msg := WSMessage{
		Type: msgType,
		Data: dataBytes,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ws broadcastAll marshal error: %v", err)
		return
	}

	h.mu.RLock()
	targets := make([]*Client, 0)
	for _, clients := range h.rooms {
		for c := range clients {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()

	for _, client := range targets {
		select {
		case client.send <- msgBytes:
		default:
			go func(c *Client) {
				h.unregister <- c
			}(client)
		}
	}
}

// ClientCount 获取房间客户端数量
func (h *WSHub) ClientCount(gameID uint64) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[gameID])
}

// readPump 读取客户端消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// MVP: 暂不处理客户端消息（subscribe/unsubscribe 通过 URL 参数实现）
	}
}

// writePump 向客户端发送消息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
