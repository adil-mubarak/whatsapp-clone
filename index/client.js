const token = localStorage.getItem('jwtToken'); // Ensure token is stored

const clientID = prompt("Enter your client ID:");
const targetID = prompt("Enter target client ID:");

const socket = new WebSocket(`ws://localhost:8081/signal?client_id=${clientID}&target_id=${targetID}`);

let localStream, peerConnection, remoteStream;

const startCallButton = document.getElementById("startCall");
const hangupCallButton = document.getElementById("hangupCall");
const localVideo = document.getElementById("localVideo");
const remoteVideo = document.getElementById("remoteVideo");

// WebSocket connection opened
socket.onopen = () => {
    console.log("WebSocket connection established.");

    // Send JWT token as a JSON message for authentication
    socket.send(JSON.stringify({ type: 'authenticate', token }));
};

// Handle messages from the server (offer, answer, candidate)
socket.onmessage = async (event) => {
    const signal = JSON.parse(event.data);
    if (signal.type === "offer") {
        await handleOffer(signal);
    } else if (signal.type === "answer") {
        await handleAnswer(signal);
    } else if (signal.type === "candidate") {
        await handleCandidate(signal);
    }
};

// Get user's media (video and audio)
async function getUserMedia() {
    try {
        localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        localVideo.srcObject = localStream;
    } catch (err) {
        console.error("Error accessing media devices:", err);
    }
}

// Create PeerConnection for WebRTC
function createPeerConnection() {
    peerConnection = new RTCPeerConnection({
        iceServers: [{ urls: "stun:stun.l.google.com:19302" }]
    });

    localStream.getTracks().forEach(track => peerConnection.addTrack(track, localStream));

    peerConnection.ontrack = (event) => {
        remoteStream = event.streams[0];
        remoteVideo.srcObject = remoteStream;
    };

    peerConnection.onicecandidate = (event) => {
        if (event.candidate) {
            sendSignal({ type: "candidate", candidate: event.candidate });
        }
    };
}

// Send signaling message to the server
function sendSignal(signal) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(signal));
    } else {
        console.error("WebSocket is not open. ReadyState is:", socket.readyState);
    }
}

// Handle Offer
async function handleOffer(signal) {
    createPeerConnection();
    await peerConnection.setRemoteDescription(new RTCSessionDescription(signal.offer));

    const answer = await peerConnection.createAnswer();
    await peerConnection.setLocalDescription(answer);

    sendSignal({ type: "answer", answer });
}

// Handle Answer
async function handleAnswer(signal) {
    const answer = new RTCSessionDescription(signal.answer);
    await peerConnection.setRemoteDescription(answer);
}

// Handle ICE Candidate
async function handleCandidate(signal) {
    const candidate = new RTCIceCandidate(signal.candidate);
    await peerConnection.addIceCandidate(candidate);
}

// Start the call
startCallButton.addEventListener("click", async () => {
    await getUserMedia();
    createPeerConnection();

    const offer = await peerConnection.createOffer();
    await peerConnection.setLocalDescription(offer);

    sendSignal({ type: "offer", offer });

    startCallButton.style.display = "none";
    hangupCallButton.style.display = "inline-block";
});

// Hang up the call
hangupCallButton.addEventListener("click", () => {
    if (peerConnection) {
        peerConnection.close();
        peerConnection = null;
    }
    remoteVideo.srcObject = null;
    startCallButton.style.display = "inline-block";
    hangupCallButton.style.display = "none";
});