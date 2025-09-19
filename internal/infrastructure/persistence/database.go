package persistence

import (
	"fmt"
	"os"
	"strconv"

	"haoma/internal/config"
	"haoma/internal/domain/leaderboard"
	"haoma/internal/domain/player"
	"haoma/internal/domain/question"
	"haoma/internal/domain/session"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase() (*Database, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", strconv.Itoa(config.DEFAULT_DB_PORT))
	dbname := getEnv("DB_NAME", "haoma")
	user := getEnv("DB_USER", "haoma")
	password := getEnv("DB_PASSWORD", "haoma_secret")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	err = db.AutoMigrate(
		&session.Session{},
		&question.Question{},
		&question.Category{},
		&player.Player{},
		&player.Attempt{},
		&leaderboard.Entry{},
	)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetJWTSecret() string {
	return getEnv("JWT_SECRET", "super_secret_key")
}
