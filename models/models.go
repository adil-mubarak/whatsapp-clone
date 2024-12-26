package models

import (
	"time"
)

type User struct {
	ID              uint   `gorm:"primaryKey"`
	PhoneNumber     string `gorm:"unique;not null" json:"phone_number"`
	UserName        string `gorm:"not null;varchar(50)" json:"user_name"`
	Profile_Picture string `gorm:"not null" json:"profile_picture"`
	OTP             string `grom:"not null" json:"otp"`
	OTPExpiry       time.Time
	CreatedAT       time.Time
	UpdateAt        time.Time
}

type Message struct {
	ID         uint      `gorm:"primarykey;autoIncrement" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	Sender     User      `gorm:"foreignKey:SenderID" json:"sender"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Receiver   User      `gorm:"foreignKey:ReceiverID" json:"receiver"`
	GroupID    *uint64   `gorm:"index;null" json:"group_id,omitempty"`
	Group      Group     `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Content    string    `gorm:"type:text" json:"content"`
	MediaURL   string    `gorm:"size:255" json:"media_url,omitempty"`
	Timestamp  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"timestamp"`
}

	type Group struct {
		ID          uint64    `gorm:"primarykey;autoIncrement;" json:"id"`  
		Name        string    `gorm:"size:100;not null" json:"name"`
		Description string    `gorm:"type:text" json:"description"`
		ProfileURL  string    `gorm:"size:255" json:"profile_url"`
		Status      string    `gorm:"type:text" json:"status"`
		AdminID     uint      `gorm:"not null" json:"admin_id"`
		Admin       User      `gorm:"foreignKey:AdminID" json:"admin"`
		CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"created_at"`
		UpdatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP;not null" json:"updated_at"`
	}


	type GroupMember struct {
		ID       uint64    `gorm:"primarykey;autoIncrement;type:bigint unsigned" json:"id"`
		GroupID  uint64    `gorm:"not null;type:bigint unsigned" json:"group_id"`
		Group    Group     `gorm:"foreignKey:GroupID" json:"group"`
		UserID   uint      `gorm:"not null" json:"user_id"`
		User     User      `gorm:"foreignKey:UserID" json:"user"`
		IsAdmin  bool      `gorm:"default:false" json:"is_admin"`
		JoinedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null" json:"joined_at"`
	}

type MessageStatus struct {
	ID        uint      `gorm:"primaryKey"`
	MessageID uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	Status    string    `gorm:"type:enum('sent', 'delivered', 'read');default:'sent'"`
	Timestamp time.Time `gorm:"autoCreateTime"`
	Message   Message   `gorm:"foreignKey:MessageID"`
	User      User      `gorm:"foreignKey:UserID"`
}
type Contact struct {
	ID          uint      `gorm:"primaryKey;autoIncrement:false;not null"`
	PhoneNumber string    `gorm:"uinque;not null" json:"phone_number"`
	AddedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	User        User      `gorm:"foreignKey:UserID"`
	Contact     User      `gorm:"foreignKey:ContactID"`
}

type UserActivity struct {
	UserID   uint      `gorm:"primaryKey;autoIncrement:false;not null"`
	LastSeen time.Time `gorm:"autoCreateTime"`
	User     User      `gorm:"foreignKey:UserID"`
}

type StatusUpdate struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"type:text" json:"user"`
	MediaURL  string    `gorm:"size:255" json:"media_url"`
	ExpiresAt time.Time `gorm:"type:timestamp" json:"expires_at"`
}
type BlockedUser struct {
	UserID        uint      `gorm:"primaryKey;autoIncrement:false"`
	BlockedUserID uint      `gorm:"primaryKey;autoIncrement:false"`
	BlockedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	User          User      `gorm:"foreignKey:UserID"`
	BlockedUser   User      `gorm:"foreignKey:BlockedUserID"`
}

// type Notification struct {
// 	ID         uint      `gorm:"primaryKey"`
// 	UserID     uint      `gorm:"not null"`
// 	Message    string    `gorm:"type:text;not null"`
// 	Type       string    `gorm:"type:enum('message', 'group', 'status');default:'message'"`
// 	ReadStatus bool      `gorm:"default:false"`
// 	Timestamp  time.Time `gorm:"autoCreateTime"`
// 	User       User      `gorm:"foreignKey:UserID"`
// }

type Setting struct {
	UserID              uint      `gorm:"primaryKey;autoIncrement:false;not null"`
	OnlineStatus        bool      `gorm:"default:true"`
	NotificationEnabled bool      `gorm:"default:true"`
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
}
