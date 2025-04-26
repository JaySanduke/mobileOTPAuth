package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateOTP(length int) string {
	digits := "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}
