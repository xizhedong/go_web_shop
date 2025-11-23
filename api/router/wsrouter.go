package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clients  = make(map[*websocket.Conn]bool)
	mu       sync.Mutex
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 消息
type Message struct {
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}

func broadcastMessage(message Message) {
	data, _ := json.Marshal(message)
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	log.Printf("new client connected: %v", conn.RemoteAddr())
	mu.Unlock()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}
		if messageType == websocket.TextMessage {
			log.Printf("received message: %s", message)
			//从控制台输入然后发送消息
			fmt.Printf("请输入要发送的消息：")
			var input string
			fmt.Scanln(&input)
			msg := Message{
				Type:      "chat",
				Content:   input,
				User:      "admin",
				Timestamp: time.Now(),
			}
			broadcastMessage(msg)
		}
	}
}

func WsRouter(r *gin.Engine) {
	ws := r.Group("/ws")
	log.Printf("ws router")
	{
		ws.GET("/", func(c *gin.Context) {
			handleChat(c.Writer, c.Request)
		})
	}
}
