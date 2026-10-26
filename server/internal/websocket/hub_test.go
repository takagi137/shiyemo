package websocket

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// Message 测试
// ============================================================

func TestMessage_Serialization(t *testing.T) {
	tests := []struct {
		name    string
		message Message
	}{
		{
			name: "带UserID的消息",
			message: Message{
				Type:    "notification",
				UserID:  123,
				Content: map[string]string{"text": "hello"},
			},
		},
		{
			name: "广播消息",
			message: Message{
				Type:    "broadcast",
				Content: "test content",
			},
		},
		{
			name: "空内容消息",
			message: Message{
				Type:    "ping",
				UserID:  456,
				Content: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 序列化
			data, err := json.Marshal(tt.message)
			require.NoError(t, err)
			assert.NotEmpty(t, data)

			// 反序列化
			var decoded Message
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err)

			assert.Equal(t, tt.message.Type, decoded.Type)
			assert.Equal(t, tt.message.UserID, decoded.UserID)
		})
	}
}

func TestMessage_UserIDOmitempty(t *testing.T) {
	msg := Message{
		Type:    "broadcast",
		Content: "test",
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	// UserID 为 0 时应该被省略
	assert.NotContains(t, string(data), "user_id")

	var decoded map[string]interface{}
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	_, exists := decoded["user_id"]
	assert.False(t, exists)
}

// ============================================================
// Hub 基础测试
// ============================================================

func TestHub_Initialization(t *testing.T) {
	hub := &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}

	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
	assert.NotNil(t, hub.broadcast)
	assert.Equal(t, 0, hub.GetOnlineCount())
}

func TestHub_GetOnlineCount(t *testing.T) {
	hub := &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}

	// 初始状态
	assert.Equal(t, 0, hub.GetOnlineCount())

	// 模拟添加客户端
	hub.mu.Lock()
	hub.clients[1] = []*Client{{userID: 1}}
	hub.clients[2] = []*Client{{userID: 2}, {userID: 2}}
	hub.mu.Unlock()

	// 应该有 2 个用户在线（不是 3 个连接）
	assert.Equal(t, 2, hub.GetOnlineCount())
}

func TestHub_IsUserOnline(t *testing.T) {
	hub := &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}

	// 用户不在线
	assert.False(t, hub.IsUserOnline(1))

	// 添加用户
	hub.mu.Lock()
	hub.clients[1] = []*Client{{userID: 1}}
	hub.mu.Unlock()

	// 用户在线
	assert.True(t, hub.IsUserOnline(1))
	assert.False(t, hub.IsUserOnline(2))
}

func TestHub_IsUserOnline_EmptyClients(t *testing.T) {
	hub := &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}

	// 添加空客户端列表
	hub.mu.Lock()
	hub.clients[1] = []*Client{}
	hub.mu.Unlock()

	// 客户端列表为空，应该返回 false
	assert.False(t, hub.IsUserOnline(1))
}

// ============================================================
// Hub 并发安全测试
// ============================================================

func TestHub_ConcurrentAccess(t *testing.T) {
	hub := &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}

	done := make(chan bool, 100)

	// 并发读写
	for i := 0; i < 50; i++ {
		go func(userID int64) {
			hub.mu.Lock()
			hub.clients[userID] = []*Client{{userID: userID}}
			hub.mu.Unlock()
			done <- true
		}(int64(i))
	}

	for i := 0; i < 50; i++ {
		go func(userID int64) {
			_ = hub.IsUserOnline(userID)
			_ = hub.GetOnlineCount()
			done <- true
		}(int64(i))
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 100; i++ {
		<-done
	}

	// 验证没有 panic
	assert.True(t, true)
}

// ============================================================
// Client 测试
// ============================================================

func TestClient_Creation(t *testing.T) {
	hub := &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}

	client := &Client{
		hub:    hub,
		userID: 123,
		send:   make(chan []byte, 256),
	}

	assert.Equal(t, hub, client.hub)
	assert.Equal(t, int64(123), client.userID)
	assert.NotNil(t, client.send)
}

func TestClient_SendChannel(t *testing.T) {
	client := &Client{
		send: make(chan []byte, 256),
	}

	// 发送消息
	testData := []byte(`{"type":"test"}`)
	client.send <- testData

	// 接收消息
	received := <-client.send
	assert.Equal(t, testData, received)
}

func TestClient_SendChannelFull(t *testing.T) {
	client := &Client{
		send: make(chan []byte, 2), // 小缓冲区
	}

	// 填满缓冲区
	client.send <- []byte("msg1")
	client.send <- []byte("msg2")

	// 非阻塞发送应该失败
	select {
	case client.send <- []byte("msg3"):
		t.Error("should not send when channel is full")
	default:
		// 预期行为
	}
}

// ============================================================
// Hub SendToUser 测试
// ============================================================

func TestHub_SendToUser(t *testing.T) {
	hub := &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}

	// 启动 Hub
	go hub.run()

	// 创建客户端
	client := &Client{
		hub:    hub,
		userID: 123,
		send:   make(chan []byte, 256),
	}

	// 注册客户端
	hub.register <- client
	time.Sleep(50 * time.Millisecond) // 等待注册完成

	// 发送消息
	hub.SendToUser(123, "test", "hello")

	// 接收消息
	select {
	case data := <-client.send:
		var msg Message
		err := json.Unmarshal(data, &msg)
		require.NoError(t, err)
		assert.Equal(t, "test", msg.Type)
		assert.Equal(t, int64(123), msg.UserID)
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for message")
	}
}

func TestHub_Broadcast(t *testing.T) {
	hub := &Hub{
		clients:    make(map[int64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}

	// 启动 Hub
	go hub.run()

	// 创建多个客户端
	client1 := &Client{
		hub:    hub,
		userID: 1,
		send:   make(chan []byte, 256),
	}
	client2 := &Client{
		hub:    hub,
		userID: 2,
		send:   make(chan []byte, 256),
	}

	// 注册客户端
	hub.register <- client1
	hub.register <- client2
	time.Sleep(50 * time.Millisecond) // 等待注册完成

	// 广播消息
	hub.Broadcast("announcement", "test broadcast")

	// 两个客户端都应该收到消息
	for _, client := range []*Client{client1, client2} {
		select {
		case data := <-client.send:
			var msg Message
			err := json.Unmarshal(data, &msg)
			require.NoError(t, err)
			assert.Equal(t, "announcement", msg.Type)
		case <-time.After(1 * time.Second):
			t.Error("timeout waiting for broadcast message")
		}
	}
}

// ============================================================
// Init 函数测试
// ============================================================

func TestInit_GlobalHub(t *testing.T) {
	// 重置全局 Hub
	globalHub = nil

	// 初始化
	Init()

	// 验证全局 Hub 已创建
	hub := GetHub()
	assert.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
	assert.NotNil(t, hub.broadcast)
}

func TestGetHub_ReturnsSameInstance(t *testing.T) {
	// 重置全局 Hub
	globalHub = nil

	Init()

	hub1 := GetHub()
	hub2 := GetHub()

	// 应该返回同一个实例
	assert.Equal(t, hub1, hub2)
}
