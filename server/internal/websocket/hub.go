package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub WebSocket 连接管理中心
type Hub struct {
	// 用户连接映射：userID -> []*Client
	clients map[int64][]*Client

	// 注册通道
	register chan *Client

	// 注销通道
	unregister chan *Client

	// 广播通道
	broadcast chan *Message

	// 互斥锁
	mu sync.RWMutex
}

// Client WebSocket 客户端
type Client struct {
	hub    *Hub
	userID int64
	conn   *websocket.Conn
	send   chan []byte
}

// Message WebSocket 消息
type Message struct {
	Type    string      `json:"type"`
	UserID  int64       `json:"user_id,omitempty"`
	Content interface{} `json:"content"`
}

// 全局 Hub 实例
var globalHub *Hub

// Init 初始化 WebSocket Hub
func Init() {
	globalHub = &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}
	go globalHub.run()
}

// GetHub 获取全局 Hub
func GetHub() *Hub {
	return globalHub
}

// run 运行 Hub
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = append(h.clients[client.userID], client)
			h.mu.Unlock()
			log.Printf("WebSocket client connected: user %d", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				for i, c := range clients {
					if c == client {
						h.clients[client.userID] = append(clients[:i], clients[i+1:]...)
						break
					}
				}
				if len(h.clients[client.userID]) == 0 {
					delete(h.clients, client.userID)
				}
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected: user %d", client.userID)

		case message := <-h.broadcast:
			h.mu.RLock()
			if message.UserID > 0 {
				// 发送给指定用户
				if clients, ok := h.clients[message.UserID]; ok {
					data, _ := json.Marshal(message)
					for _, client := range clients {
						select {
						case client.send <- data:
						default:
							close(client.send)
						}
					}
				}
			} else {
				// 广播给所有用户
				data, _ := json.Marshal(message)
				for _, clients := range h.clients {
					for _, client := range clients {
						select {
						case client.send <- data:
						default:
							close(client.send)
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser 发送消息给指定用户
func (h *Hub) SendToUser(userID int64, msgType string, content interface{}) {
	h.broadcast <- &Message{
		Type:    msgType,
		UserID:  userID,
		Content: content,
	}
}

// Broadcast 广播消息给所有用户
func (h *Hub) Broadcast(msgType string, content interface{}) {
	h.broadcast <- &Message{
		Type:    msgType,
		Content: content,
	}
}

// GetOnlineCount 获取在线用户数
func (h *Hub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// IsUserOnline 检查用户是否在线
// SetGlobalHubForTest 设置全局 Hub（仅测试用）
func SetGlobalHubForTest(h *Hub) {
	globalHub = h
}

// NewTestHub 创建测试用 Hub（不启动 run goroutine）
func NewTestHub() *Hub {
	return &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}
}

// AddOnlineUserForTest 直接添加在线用户（仅测试用）
func (h *Hub) AddOnlineUserForTest(userID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = append(h.clients[userID], &Client{})
}

func (h *Hub) IsUserOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}
