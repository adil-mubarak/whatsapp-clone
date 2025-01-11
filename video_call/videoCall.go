package main

import (
	"log"
	"net/http"
	"sync"
	"whatsapp/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var clients = make(map[string]*Client)
var clientsMutex = sync.Mutex{}

func signalHandler(c *gin.Context) {
	claims, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims, ok := claims.(*service.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}
	userID := userClaims.PhoneNumber
	log.Println(userID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data"})
		return
	}

	// var userphone models.User
	// if err := db.DB.Where("id = ?",userID).Find(&userphone).Error; err != nil{
	// 	c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
	// 	return
	// }
	// clientID := userphone.PhoneNumber
	clientID := c.Query("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	client := &Client{ID: clientID, Conn: conn}
	clientsMutex.Lock()
	clients[clientID] = client
	clientsMutex.Unlock()
	log.Printf("Client connected: %s", clientID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", clientID, err)
			break
		}

		log.Printf("Received from %s: %s", clientID, string(msg))

		targetID := c.Query("target_id")
		if targetID == "" {
			log.Println("target_id is required in message")
			continue
		}

		clientsMutex.Lock()
		targetClient, exists := clients[targetID]
		clientsMutex.Unlock()
		if exists {
			err = targetClient.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error sending message to %s: %v", targetID, err)
			}
		} else {
			log.Printf("Target client %s not found", targetID)
		}
	}

	clientsMutex.Lock()
	delete(clients, clientID)
	clientsMutex.Unlock()
	log.Printf("Client disconnected: %s", clientID)
}

func main() {
	r := gin.Default()

	r.GET("/signal", signalHandler)

	log.Println("WebSocket server running on ws://localhost:8080/signal...")
	log.Fatal(r.Run(":8081"))
}
