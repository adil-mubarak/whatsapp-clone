package service

import (
	"log"
	"net/http"
	"time"
	"whatsapp/db"
	"whatsapp/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)

func WebSocketHandler(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error during WebSocket upgrade:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not open WebSocket connection"})
		return
	}

	clients[conn] = true
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			delete(clients, conn)
			break
		}

		msg.Timestamp = time.Now()

		if err := db.DB.Create(&msg).Error; err != nil {
			log.Println("Error saving message to the database:", err)
			continue
		}

		broadcast <- msg
	}
}

func BroadCastMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error sending message to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
