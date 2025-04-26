package service

import (
	"context"
	"fmt"
)

type SMSService interface {
	SendOTP(ctx context.Context, mobile string, otp string) error
}

type smsService struct{}

func NewSMSService() SMSService {
	return &smsService{}
}

func (s *smsService) SendOTP(ctx context.Context, mobileNumber string, otp string) error {
	fmt.Printf("Dummy SMS service, Sent OTP %s to the mobile number %s\n", otp, mobileNumber)
	return nil
}
