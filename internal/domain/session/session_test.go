package session

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSession_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		session  Session
		expected bool
	}{
		{
			name: "active session within time limit",
			session: Session{
				ID:        uuid.New(),
				StartedAt: time.Now().Add(-30 * time.Minute),
			},
			expected: true,
		},
		{
			name: "expired session beyond 2 hours",
			session: Session{
				ID:        uuid.New(),
				StartedAt: time.Now().Add(-3 * time.Hour),
			},
			expected: false,
		},
		{
			name: "finished session",
			session: Session{
				ID:         uuid.New(),
				StartedAt:  time.Now().Add(-1 * time.Hour),
				FinishedAt: &[]time.Time{time.Now().Add(-30 * time.Minute)}[0],
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.IsActive(); got != tt.expected {
				t.Errorf("Session.IsActive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSession_CalculateScore(t *testing.T) {
	startTime := time.Now().Add(-90 * time.Second) // 90 seconds ago

	s := &Session{
		ID:        uuid.New(),
		StartedAt: startTime,
		Score: Score{
			Correct: 5,
			Total:   7,
		},
	}

	score := s.CalculateScore()

	expectedTimePenalty := 3                         // 90 seconds / 30 = 3 points
	expectedFinal := (5 * 100) - expectedTimePenalty // 500 - 3 = 497

	if score.TimePenalty != expectedTimePenalty {
		t.Errorf("Expected time penalty %d, got %d", expectedTimePenalty, score.TimePenalty)
	}

	if score.Final != expectedFinal {
		t.Errorf("Expected final score %d, got %d", expectedFinal, score.Final)
	}

	if score.Correct != 5 {
		t.Errorf("Expected correct answers 5, got %d", score.Correct)
	}
}
