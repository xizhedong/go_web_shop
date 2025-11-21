// 我想要建立一个websocket服务器，使用Go语言编写，能够处理客户端的连接请求，并且能够接收和发送消息。请帮我完成这个代码。
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var(
	clients =make(map[*websocket.Conn]bool)
	mu sync.Mutex
	upgrader =websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func broadcastMessage(message []byte) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
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
			broadcastMessage(message)
		}
	}
}	



func main() {
	http.HandleFunc("/ws", handleChat)
	go handleChat()
	log.Println("服务端启动：ws://localhost:8080/ws")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("服务启动失败：%v", err)
	}
}
