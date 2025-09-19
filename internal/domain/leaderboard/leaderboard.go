package leaderboard

import (
	"time"

	"github.com/google/uuid"
)

// Entry represents a champion's achievement on the taxteh-ye sharaf
type Entry struct {
	ID             uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	PlayerID       uuid.UUID     `json:"player_id" gorm:"type:uuid;not null"`
	PlayerName     string        `json:"player_name" gorm:"not null"`
	SessionID      uuid.UUID     `json:"session_id" gorm:"type:uuid;not null"`
	FinalScore     int           `json:"final_score" gorm:"not null"`
	CompletionTime time.Duration `json:"completion_time" gorm:"type:bigint"` // For tie-breaking
	AchievedAt     time.Time     `json:"achieved_at"`
}

// Leaderboard maintains the eternal witness of glory
type Leaderboard struct {
	Entries []Entry `json:"entries"`
}

func (l *Leaderboard) Top10() []Entry {
	if len(l.Entries) <= 10 {
		return l.Entries
	}
	return l.Entries[:10]
}

func NewEntry(playerID uuid.UUID, playerName string, sessionID uuid.UUID,
	finalScore int, completionTime time.Duration) *Entry {
	return &Entry{
		ID:             uuid.New(),
		PlayerID:       playerID,
		PlayerName:     playerName,
		SessionID:      sessionID,
		FinalScore:     finalScore,
		CompletionTime: completionTime,
		AchievedAt:     time.Now(),
	}
}
