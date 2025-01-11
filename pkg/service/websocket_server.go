package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
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
var clientsMutex sync.Mutex

func handleSignal(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer conn.Close()

	client := &Client{ID: clientID, Conn: conn}

	clientsMutex.Lock()
	clients[clientID] = client
	clientsMutex.Unlock()

	log.Printf("Client connected: %s", clientID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client %s: %v", clientID, err)
			break
		}

		var signal map[string]interface{}
		if err := json.Unmarshal(message, &signal); err != nil {
			log.Printf("Error decoding JSON message: %v", err)
			continue
		}

		targetID, ok := signal["target_id"].(string)
		if !ok || targetID == "" {
			log.Printf("target_id is missing or invalid")
			continue
		}

		clientsMutex.Lock()
		targetClient, exists := clients[targetID]
		clientsMutex.Unlock()

		if exists {
			err := targetClient.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error sending message to client %s: %v", targetID, err)
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

func handleWebRTC(w http.ResponseWriter, r *http.Request) {
	var offer webrtc.SessionDescription
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		http.Error(w, "Invalid offer", http.StatusBadRequest)
		return
	}

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		log.Printf("Error creating peer connection: %v", err)
		http.Error(w, "Failed to create peer connection", http.StatusInternalServerError)
		return
	}

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			log.Printf("New ICE candidate: %s", candidate.ToJSON().Candidate)
		}
	})

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Received track: %s", track.Kind().String())
	})

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Printf("Error setting remote description: %v", err)
		http.Error(w, "Failed to set remote description", http.StatusInternalServerError)
		return
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Printf("Error creating answer: %v", err)
		http.Error(w, "Failed to create answer", http.StatusInternalServerError)
		return
	}

	if err := peerConnection.SetLocalDescription(answer); err != nil {
		log.Printf("Error setting local description: %v", err)
		http.Error(w, "Failed to set local description", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		log.Printf("Error encoding answer: %v", err)
		http.Error(w, "Failed to encode answer", http.StatusInternalServerError)
		return
	}
}

// func main() {
// 	http.HandleFunc("/signal", handleSignal)
// 	http.HandleFunc("/webrtc", handleWebRTC)

// 	log.Println("Server is running on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
