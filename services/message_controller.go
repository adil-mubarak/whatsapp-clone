package services

import (
	"context"
	"log"
	"net/http"
	"time"
	"whatsapp/db"
	"whatsapp/models"
	proto "whatsapp/proto/chatproto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChatServiceServer struct {
	DB *gorm.DB
	proto.UnimplementedChatServiceServer
}

func SendMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.Timestamp = time.Now()

	if err := db.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully", "data": message})
}

func GetMessage(c *gin.Context) {
	senderID := c.Query("sender_id")
	receiverID := c.Query("receiver_id")

	if senderID == "" || receiverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender and Receiver IDs are required"})
		return
	}

	var messages []models.Message
	if err := db.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID).Order("timestamp asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}


func (s *ChatServiceServer)SendMessage(ctx context.Context,req *proto.MessageRequest)(*proto.MessageResponse,error){

	message := models.Message{
		SenderID: uint(req.SenderId),
		ReceiverID: uint(req.ReceiverId),
		GroupID: nil,
		Content: req.Content,
		MediaURL: req.MediaUrl,
		Timestamp: time.Now(),
	}

	if req.GroupId != 0 {
		log.Printf("Grouped provided: %d",req.GroupId)
		groupID := uint64(req.GroupId)
		message.GroupID = &groupID
	}

	if err := s.DB.Create(&message).Error; err != nil{
		return &proto.MessageResponse{
			Success: false,
			Message: "Failed to save message: "+err.Error(),
		},nil
	}

	return &proto.MessageResponse{
		Success: true,
		Message: "Message sent successfully",
	},nil
}

func (s *ChatServiceServer) ChatStream(stream proto.ChatService_ChatStreamServer)error{
	for{
		incomingMessage, err := stream.Recv()
		if err != nil{
			log.Printf("Error receiving message: %v",err)
			return err
		}

		log.Printf("Received message: %+v",incomingMessage)

		message := models.Message{
			SenderID: uint(incomingMessage.SenderId),
			ReceiverID: uint(incomingMessage.ReceiverId),
			GroupID: nil,
			Content: incomingMessage.Content,
			MediaURL: incomingMessage.MediaUrl,
			Timestamp: time.Now(),
		}

		if err := s.DB.Create(&message).Error; err != nil{
			log.Printf("Failed to save message: %v",err)
			return err
		}

		if err := stream.Send(&proto.ChatMessage{
			SenderId: incomingMessage.SenderId,
			ReceiverId: incomingMessage.ReceiverId,
			GroupId: incomingMessage.GroupId,
			Content: incomingMessage.Content,
			MediaUrl: incomingMessage.MediaUrl,
			Timestamp: message.Timestamp.Format(time.RFC3339),
		});err != nil{
			return err
		}
	}
}