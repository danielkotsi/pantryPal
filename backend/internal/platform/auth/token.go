package auth

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
	UserID string `json:"uid"`
	Exp    int64  `json:"exp"`
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl}
}

func (t *TokenManager) MakeToken(userID string) (string, time.Time, error) {
	expiresAt := time.Now().UTC().Add(t.ttl)
	claims := Claims{UserID: userID, Exp: expiresAt.Unix()}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	payload64 := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, t.secret)
	if _, err := mac.Write([]byte(payload64)); err != nil {
		return "", time.Time{}, err
	}
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return payload64 + "." + sig, expiresAt, nil
}

func (t *TokenManager) ParseToken(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Claims{}, errors.New("invalid token format")
	}

	mac := hmac.New(sha256.New, t.secret)
	if _, err := mac.Write([]byte(parts[0])); err != nil {
		return Claims{}, err
	}
	expected := mac.Sum(nil)
	actual, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, err
	}
	if !hmac.Equal(expected, actual) {
		return Claims{}, errors.New("token signature mismatch")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, err
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, err
	}
	if claims.UserID == "" || claims.Exp == 0 {
		return Claims{}, errors.New("invalid token claims")
	}
	if time.Now().UTC().Unix() > claims.Exp {
		return Claims{}, errors.New("token expired")
	}
	return claims, nil
}
