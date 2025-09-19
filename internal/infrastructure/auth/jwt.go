package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey []byte
}

type PlayerClaims struct {
	PlayerID    uuid.UUID  `json:"player_id"`
	PlayerName  string     `json:"player_name"`
	PlayerEmail string     `json:"player_email"`
	SessionID   *uuid.UUID `json:"session_id,omitempty"`
	jwt.RegisteredClaims
}

type SessionClaims struct {
	PlayerID  uuid.UUID `json:"player_id"`
	SessionID uuid.UUID `json:"session_id"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

func (j *JWTService) GeneratePlayerToken(playerID uuid.UUID, playerName, playerEmail string) (string, error) {
	claims := PlayerClaims{
		PlayerID:    playerID,
		PlayerName:  playerName,
		PlayerEmail: playerEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hour expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "haoma-carnival",
			Subject:   playerID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) GenerateSessionToken(playerID, sessionID uuid.UUID) (string, error) {
	claims := SessionClaims{
		PlayerID:  playerID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 2 hour session limit
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "haoma-carnival",
			Subject:   sessionID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidatePlayerToken(tokenString string) (*PlayerClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PlayerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*PlayerClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWTService) ValidateSessionToken(tokenString string) (*SessionClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SessionClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid session token")
}

func (j *JWTService) RefreshPlayerToken(oldTokenString string) (string, error) {
	claims, err := j.ValidatePlayerToken(oldTokenString)
	if err != nil {
		return "", err
	}

	return j.GeneratePlayerToken(claims.PlayerID, claims.PlayerName, claims.PlayerEmail)
}
