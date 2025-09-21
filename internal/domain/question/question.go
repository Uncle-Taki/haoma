package question

import (
	"log"

	"github.com/google/uuid"
)

// Question represents riddles with answers
type Question struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Text        string    `json:"text" gorm:"type:text;not null"`
	OptionA     string    `json:"option_a" gorm:"type:text;not null"`
	OptionB     string    `json:"option_b" gorm:"type:text;not null"`
	OptionC     *string   `json:"option_c,omitempty" gorm:"type:text"`
	OptionD     *string   `json:"option_d,omitempty" gorm:"type:text"`
	Correct     string    `json:"-" gorm:"not null"`
	Explanation string    `json:"explanation" gorm:"type:text"`
	CategoryID  uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
}

// Category represents knowledge domains
type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name        string    `json:"name" gorm:"type:text;unique;not null"`
	Title       string    `json:"title" gorm:"type:text"`
	Description string    `json:"description" gorm:"type:text"`
	IsPhDT      bool      `json:"is_phdt" gorm:"default:false"`
}

// Node represents a trial tent holding questions
type Node struct {
	Number     int        `json:"number"`
	CategoryID uuid.UUID  `json:"category_id"`
	Questions  []Question `json:"questions"`
}

func (question *Question) IsBinaryChoice() bool {
	return question.OptionC == nil && question.OptionD == nil
}

func (question *Question) ValidateAnswer(answer string) bool {
	log.Println("Correct answer:", question.Correct)
	log.Println("Answer:", answer)
	return question.Correct == answer
}
