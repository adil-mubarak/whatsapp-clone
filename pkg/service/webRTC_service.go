package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pion/webrtc/v3"
)

var peerConnection = make(map[string]*webrtc.PeerConnection)

func HandleOffer(w http.ResponseWriter, r *http.Request) {
	var offer webrtc.SessionDescription
	err := json.NewDecoder(r.Body).Decode(&offer)
	if err != nil {
		http.Error(w, "Invalid offer", http.StatusBadRequest)
		return
	}

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			log.Printf("New ICE candidaate: %s", candidate.ToJSON().Candidate)
		}
	})

	err = pc.SetRemoteDescription(offer)
	if err != nil {
		log.Fatal(err)
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = pc.SetLocalDescription(answer)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Contetn-Type", "application/json")
	json.NewEncoder(w).Encode(answer)

	peerConnection[pc.LocalDescription().SDP] = pc
}
