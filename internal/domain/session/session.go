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
func (slice StringSlice) Value() (driver.Value, error) {
	if slice == nil {
		return nil, nil
	}
	return json.Marshal([]string(slice))
}

// Scan implements sql.Scanner interface for database deserialization
func (slice *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*slice = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, (*[]string)(slice))
}

// IntMap is a custom type for proper JSON serialization with PostgreSQL
type IntMap map[int]int64

// Value implements driver.Valuer interface for database serialization
func (intMap IntMap) Value() (driver.Value, error) {
	if intMap == nil {
		return nil, nil
	}
	return json.Marshal(map[int]int64(intMap))
}

// Scan implements sql.Scanner interface for database deserialization
func (intMap *IntMap) Scan(value interface{}) error {
	if value == nil {
		*intMap = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, (*map[int]int64)(intMap))
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

func (session *Session) IsActive() bool {
	if session.FinishedAt != nil {
		return false
	}

	elapsed := time.Since(session.StartedAt)
	return elapsed < config.MAX_SESSION_DURATION
}

func (session *Session) CalculateScore() Score {
	// Use the accumulated time penalty from completed nodes
	final := (session.Score.Correct * config.CORRECT_ANSWER_MULTIPLIER) -
		(session.Score.TimePenalty * config.PENALTY_MULTIPLIER)
	if final < 0 {
		final = config.DEFAULT_SCORE
	}

	return Score{
		Correct:     session.Score.Correct,
		Total:       session.Score.Total,
		TimePenalty: session.Score.TimePenalty,
		Final:       final,
	}
}
