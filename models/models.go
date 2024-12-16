package models

import "time"

type User struct {
	ID          uint   `gorm:"primaryKey"`
	PhoneNumber string `gorm:"unique;not null" json:"phone_number"`
	UserName string `gorm:"not null;varchar(50)" json:"user_name"`
	Profile_Picture string `gorm:"not null" json:"profile_picture"`
	OTP         string `grom:"not null" json:"otp"`
	OTPExpiry   time.Time
	CreatedAT   time.Time
	UpdateAt    time.Time
}
