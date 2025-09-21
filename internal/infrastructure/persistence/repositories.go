package persistence

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"haoma/internal/config"
	"haoma/internal/domain/leaderboard"
	"haoma/internal/domain/player"
	"haoma/internal/domain/question"
	"haoma/internal/domain/session"
)

// SessionRepository implements session persistence
type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Save(session *session.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) FindByID(id uuid.UUID) (*session.Session, error) {
	var foundSession session.Session
	err := r.db.First(&foundSession, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("session not found")
	}
	return &foundSession, err
}

func (r *SessionRepository) Update(session *session.Session) error {
	return r.db.Save(session).Error
}

// QuestionRepository implements question persistence
type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) GetCategories() ([]question.Category, error) {
	var categories []question.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *QuestionRepository) GetQuestionsByCategory(categoryID uuid.UUID, limit int) ([]question.Question, error) {
	var questions []question.Question
	err := r.db.Where("category_id = ?", categoryID).
		Preload("Category").
		Limit(limit).
		Find(&questions).Error
	return questions, err
}

func (r *QuestionRepository) GetUnusedFunQuestionsForSession(sessionID uuid.UUID, limit int) ([]question.Question, error) {
	var questions []question.Question

	err := r.db.Joins("JOIN categories ON questions.category_id = categories.id").
		Where("categories.name = ?", config.FUN_CATEGORY_NAME).
		Where("questions.id NOT IN (SELECT question_id FROM attempts WHERE session_id = ?)", sessionID).
		Preload("Category").
		Limit(limit).
		Find(&questions).Error

	return questions, err
}

func (r *QuestionRepository) FindByID(id uuid.UUID) (*question.Question, error) {
	var foundQuestion question.Question
	err := r.db.Preload("Category").First(&foundQuestion, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("question not found")
	}
	return &foundQuestion, err
}

// PlayerRepository implements player persistence
type PlayerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Save(player *player.Player) error {
	return r.db.Create(player).Error
}

func (r *PlayerRepository) FindByID(id uuid.UUID) (*player.Player, error) {
	var foundPlayer player.Player
	err := r.db.First(&foundPlayer, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("player not found")
	}
	return &foundPlayer, err
}

func (r *PlayerRepository) FindByEmail(email string) (*player.Player, error) {
	var foundPlayer player.Player
	err := r.db.Where("email = ?", email).First(&foundPlayer).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("player not found")
	}
	return &foundPlayer, err
}

func (r *PlayerRepository) SaveAttempt(attempt *player.Attempt) error {
	return r.db.Create(attempt).Error
}

func (r *PlayerRepository) GetAttemptsBySessionAndCategory(sessionID, categoryID uuid.UUID) ([]player.Attempt, error) {
	var attempts []player.Attempt
	err := r.db.Joins("JOIN questions ON attempts.question_id = questions.id").
		Where("attempts.session_id = ? AND questions.category_id = ?", sessionID, categoryID).
		Find(&attempts).Error
	return attempts, err
}

func (r *PlayerRepository) HasAnsweredQuestion(sessionID, questionID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&player.Attempt{}).
		Where("session_id = ? AND question_id = ?", sessionID, questionID).
		Count(&count).Error
	return count > 0, err
}

// LeaderboardRepository implements leaderboard persistence
type LeaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) *LeaderboardRepository {
	return &LeaderboardRepository{db: db}
}

func (r *LeaderboardRepository) AddEntry(entry *leaderboard.Entry) error {
	return r.db.Create(entry).Error
}

func (r *LeaderboardRepository) UpsertEntry(entry *leaderboard.Entry) error {
	var existingEntry leaderboard.Entry
	err := r.db.Where("session_id = ?", entry.SessionID).First(&existingEntry).Error

	if err == nil {
		existingEntry.FinalScore = entry.FinalScore
		existingEntry.CompletionTime = entry.CompletionTime
		existingEntry.AchievedAt = time.Now()
		return r.db.Save(&existingEntry).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(entry).Error
	} else {
		return err
	}
}

func (r *LeaderboardRepository) GetTop10() ([]leaderboard.Entry, error) {
	var entries []leaderboard.Entry
	err := r.db.Order("final_score DESC, completion_time ASC").
		Limit(config.LEADERBOARD_TOP_ENTRIES).
		Find(&entries).Error
	return entries, err
}
