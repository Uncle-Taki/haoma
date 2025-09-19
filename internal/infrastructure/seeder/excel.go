package seeder

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"haoma/internal/domain/question"
)

type ExcelSeeder struct {
	db *gorm.DB
}

func NewExcelSeeder(db *gorm.DB) *ExcelSeeder {
	return &ExcelSeeder{db: db}
}

// safeGetColumn safely accesses a column index in a row slice, returning empty string if index is out of range
func safeGetColumn(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return row[index]
}

func (s *ExcelSeeder) SeedFromExcel(scenariosPath, questionsPath string) error {
	log.Println("🌱 Seeding carnival with mystical knowledge...")

	// Load categories from SCENARIOS.xlsx
	if err := s.seedCategories(scenariosPath); err != nil {
		return err
	}

	// Load questions from questions.xlsx
	if err := s.seedQuestions(questionsPath); err != nil {
		return err
	}

	log.Println("✨ The carnival's wisdom has been awakened!")
	return nil
}

// seedCategories loads categories from SCENARIOS.xlsx
func (s *ExcelSeeder) seedCategories(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return errors.New("no sheets found in SCENARIOS file")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return err
	}

	if len(rows) < 2 {
		return errors.New("insufficient data in SCENARIOS file")
	}

	// Skip header row
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < 1 {
			continue // Skip incomplete rows
		}

		name := safeGetColumn(rows[i], 0)
		title := safeGetColumn(rows[i], 1)
		description := safeGetColumn(rows[i], 2)
		isPhDT := name == "PhDT"

		category := question.Category{
			ID:          uuid.New(),
			Name:        name,
			Title:       title,
			Description: description,
			IsPhDT:      isPhDT,
		}

		result := s.db.Where("name = ?", name).FirstOrCreate(&category)
		if result.Error != nil {
			log.Printf("Error seeding category %s: %v", name, result.Error)
			continue
		}

		log.Printf("📚 Added category: %s (PhDT: %v)", name, isPhDT)
	}

	return nil
}

// seedQuestions loads questions from questions.xlsx
func (s *ExcelSeeder) seedQuestions(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return errors.New("no sheets found in questions file")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return err
	}

	if len(rows) < 2 {
		return errors.New("insufficient data in questions file")
	}

	for i := 1; i < len(rows); i++ {
		if len(rows[i]) < 2 {
			continue // Skip incomplete rows
		}

		categoryName := safeGetColumn(rows[i], 0)
		text := safeGetColumn(rows[i], 1)
		optionA := safeGetColumn(rows[i], 2)
		optionB := safeGetColumn(rows[i], 3)
		optionC := safeGetColumn(rows[i], 4)
		optionD := safeGetColumn(rows[i], 5)
		correct := safeGetColumn(rows[i], 6)
		explanation := safeGetColumn(rows[i], 7)

		// Find category
		var category question.Category
		if err := s.db.Where("name = ?", categoryName).First(&category).Error; err != nil {
			log.Printf("Category %s not found, skipping question", categoryName)
			continue
		}

		// Handle PhDT questions (binary only)
		var optionCPtr, optionDPtr *string
		if !category.IsPhDT {
			optionCPtr = &optionC
			optionDPtr = &optionD
		} // For PhDT, leave C and D as nil

		q := question.Question{
			ID:          uuid.New(),
			Text:        text,
			OptionA:     optionA,
			OptionB:     optionB,
			OptionC:     optionCPtr,
			OptionD:     optionDPtr,
			Correct:     correct,
			Explanation: explanation,
			CategoryID:  category.ID,
		}

		result := s.db.Create(&q)
		if result.Error != nil {
			log.Printf("Error seeding question: %v", result.Error)
			continue
		}

		log.Printf("❓ Added question for %s (PhDT: %v)", categoryName, category.IsPhDT)
	}

	return nil
}
