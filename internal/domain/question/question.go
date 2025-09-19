package question

import (
	"log"

	"github.com/google/uuid"
)

// Question represents riddles with answers
type Question struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Text        string    `json:"text" gorm:"type:text;not null"`
	OptionA     string    `json:"option_a" gorm:"not null"`
	OptionB     string    `json:"option_b" gorm:"not null"`
	OptionC     *string   `json:"option_c,omitempty"`
	OptionD     *string   `json:"option_d,omitempty"`
	Correct     string    `json:"-" gorm:"not null"`
	Explanation string    `json:"explanation"`
	CategoryID  uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
}

// Category represents knowledge domains
type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPhDT      bool      `json:"is_phdt" gorm:"default:false"`
}

// Node represents a trial tent holding questions
type Node struct {
	Number     int        `json:"number"`
	CategoryID uuid.UUID  `json:"category_id"`
	Questions  []Question `json:"questions"`
}

func (q *Question) IsBinaryChoice() bool {
	return q.OptionC == nil && q.OptionD == nil
}

func (q *Question) ValidateAnswer(answer string) bool {
	log.Println("Correct answer:", q.Correct)
	log.Println("Answer:", answer)
	return q.Correct == answer
}
