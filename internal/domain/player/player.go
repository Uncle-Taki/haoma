package player

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Player represents the brave soul
type Player struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Attempt captures a player's answer in time
type Attempt struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	SessionID  uuid.UUID `json:"session_id" gorm:"type:uuid;not null"`
	QuestionID uuid.UUID `json:"question_id" gorm:"type:uuid;not null"`
	Answer     string    `json:"answer" gorm:"not null"`
	IsCorrect  bool      `json:"is_correct"`
	AttemptAt  time.Time `json:"attempt_at"`
}

func NewPlayer(name, email, password string) (*Player, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Player{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (player *Player) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(player.PasswordHash), []byte(password))
	return err == nil
}

func (player *Player) UpdatePassword(newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	player.PasswordHash = string(hashedPassword)
	player.UpdatedAt = time.Now()
	return nil
}

func NewAttempt(sessionID, questionID uuid.UUID, answer string, isCorrect bool) *Attempt {
	return &Attempt{
		ID:         uuid.New(),
		SessionID:  sessionID,
		QuestionID: questionID,
		Answer:     answer,
		IsCorrect:  isCorrect,
		AttemptAt:  time.Now(),
	}
}
