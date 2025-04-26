package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret string
}

type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	MobileNumber string `json:"mobile"`
	Fingerprint  string `json:"fingerprint"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{Secret: secret}
}

func (jm *JWTManager) GenerateToken(userID uint, mobileNumber, fingerprint string, ttlMinutes int) (string, error) {
	claims := JWTClaims{
		UserID:       userID,
		MobileNumber: mobileNumber,
		Fingerprint:  fingerprint,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userID), 10),
			Issuer:    "mobileOtpAuthSystem",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttlMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.Secret))
}

func (jm *JWTManager) VerifyToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jm.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
