package question

import (
	"testing"

	"github.com/google/uuid"
)

func TestQuestion_IsBinaryChoice(t *testing.T) {
	tests := []struct {
		name     string
		question Question
		expected bool
	}{
		{
			name: "PhDT question (binary)",
			question: Question{
				ID:      uuid.New(),
				Text:    "Is this phishing?",
				OptionA: "Yes",
				OptionB: "No",
				OptionC: nil,
				OptionD: nil,
			},
			expected: true,
		},
		{
			name: "Regular question (multiple choice)",
			question: Question{
				ID:      uuid.New(),
				Text:    "What is AES?",
				OptionA: "Algorithm",
				OptionB: "Protocol",
				OptionC: &[]string{"Standard"}[0],
				OptionD: &[]string{"Method"}[0],
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.question.IsBinaryChoice(); got != tt.expected {
				t.Errorf("Question.IsBinaryChoice() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestQuestion_ValidateAnswer(t *testing.T) {
	q := Question{
		ID:      uuid.New(),
		Correct: "B",
	}

	tests := []struct {
		name     string
		answer   string
		expected bool
	}{
		{"correct answer", "B", true},
		{"incorrect answer", "A", false},
		{"invalid answer", "X", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := q.ValidateAnswer(tt.answer); got != tt.expected {
				t.Errorf("Question.ValidateAnswer(%s) = %v, want %v", tt.answer, got, tt.expected)
			}
		})
	}
}
