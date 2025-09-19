package session

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"haoma/internal/config"
)

// StringSlice is a custom type for proper JSON serialization with PostgreSQL
type StringSlice []string

// Value implements driver.Valuer interface for database serialization
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal([]string(s))
}

// Scan implements sql.Scanner interface for database deserialization
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, (*[]string)(s))
}

// IntMap is a custom type for proper JSON serialization with PostgreSQL
type IntMap map[int]int64

// Value implements driver.Valuer interface for database serialization
func (m IntMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(map[int]int64(m))
}

// Scan implements sql.Scanner interface for database deserialization
func (m *IntMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, (*map[int]int64)(m))
}

// Session orchestrates a player's journey through carnival nodes
type Session struct {
	ID             uuid.UUID   `json:"id" gorm:"type:uuid;primary_key"`
	PlayerID       uuid.UUID   `json:"player_id" gorm:"type:uuid;not null"`
	StartedAt      time.Time   `json:"started_at"`
	FinishedAt     *time.Time  `json:"finished_at,omitempty"`
	CurrentNode    int         `json:"current_node" gorm:"default:1"`
	Score          Score       `json:"score" gorm:"embedded"`
	Categories     StringSlice `json:"categories" gorm:"type:json"`
	NodeStartTimes IntMap      `json:"node_start_times" gorm:"type:json"`
}

// Score represents correctness minus time's cruel tax
type Score struct {
	Correct     int `json:"correct"`
	Total       int `json:"total"`
	TimePenalty int `json:"time_penalty"`
	Final       int `json:"final"`
}

// TimeWindow enforces the maximum session duration boundary
type TimeWindow struct {
	Start time.Time
	End   time.Time
}

func (s *Session) IsActive() bool {
	if s.FinishedAt != nil {
		return false
	}

	elapsed := time.Since(s.StartedAt)
	return elapsed < config.MAX_SESSION_DURATION
}

func (s *Session) CalculateScore() Score {
	// Use the accumulated time penalty from completed nodes
	final := (s.Score.Correct * config.CORRECT_ANSWER_MULTIPLIER) -
		(s.Score.TimePenalty * config.PENALTY_MULTIPLIER)
	if final < 0 {
		final = config.DEFAULT_SCORE
	}

	return Score{
		Correct:     s.Score.Correct,
		Total:       s.Score.Total,
		TimePenalty: s.Score.TimePenalty,
		Final:       final,
	}
}
