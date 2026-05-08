package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Claims struct {
	UserID    uint64 `json:"user_id"`
	Username  string `json:"username,omitempty"`
	ExpiresAt int64  `json:"exp"`
}

func GenerateToken(userID uint64, username string, secret string, expireSeconds int64) (string, error) {
	if expireSeconds <= 0 {
		expireSeconds = 7200
	}

	headerBytes, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}

	payloadBytes, err := json.Marshal(Claims{
		UserID:    userID,
		Username:  username,
		ExpiresAt: time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix(),
	})
	if err != nil {
		return "", err
	}

	header := base64.RawURLEncoding.EncodeToString(headerBytes)
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signature := sign(header+"."+payload, secret)

	return header + "." + payload + "." + signature, nil
}

func ParseToken(token string, secret string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}

	expectedSignature := sign(parts[0]+"."+parts[1], secret)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, errors.New("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

func sign(data string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
