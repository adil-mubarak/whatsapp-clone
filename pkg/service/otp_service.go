package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendOTP(phoneNumber string) (string, error) {
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	params := &twilioApi.CreateMessageParams{}
	params.SetTo("+91" + phoneNumber)
	params.SetFrom(os.Getenv("TWILIO_PHONE"))
	params.SetBody(fmt.Sprintf("Your WhatsApp verification code is: %s \nThis code is valid for 5 minutes. Please do not share it with anyone", otp))

	_, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send OTP: %v", err)
		return "", err
	}
	return otp, nil
}
