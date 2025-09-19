package config

import "time"

// ================================
// CARNIVAL GAME LOGIC CONSTANTS
// ================================

const (
	// Core game flow
	MAX_CARNIVAL_NODES              = 7 // Total number of carnival nodes/categories
	QUESTIONS_PER_NODE              = 5 // Total questions per carnival node
	CATEGORY_QUESTIONS_PER_NODE     = 4 // Category-specific questions per node
	FUN_QUESTIONS_PER_NODE          = 1 // Fun questions per node
	MIN_REQUIRED_GENERAL_CATEGORIES = 7 // Minimum general categories needed

	// Scoring system
	CORRECT_ANSWER_MULTIPLIER     = 100 // Points per correct answer
	PENALTY_MULTIPLIER            = 10  // Points per penalty point
	TIME_PENALTY_INTERVAL_SECONDS = 30  // Seconds per penalty point

	// Node validation
	MIN_NODE_NUMBER = 1 // Minimum valid node number
	MAX_NODE_NUMBER = 7 // Maximum valid node number

	// Question requirements
	QUESTIONS_TO_COMPLETE_NODE = 5 // Questions needed to complete a node
	NODES_TO_COMPLETE_SESSION  = 7 // Nodes needed to complete session

	// Database limits
	LEADERBOARD_TOP_ENTRIES   = 10 // Number of top entries in leaderboard
	QUESTION_FETCH_MULTIPLIER = 2  // Multiplier for fetching extra questions
)

// ================================
// TIME-RELATED CONSTANTS
// ================================

const (
	// Session duration
	MAX_SESSION_DURATION_HOURS = 2                                      // Maximum session duration in hours
	MAX_SESSION_DURATION       = MAX_SESSION_DURATION_HOURS * time.Hour // Maximum session duration
	SESSION_EXPIRY_SECONDS     = 7200                                   // Session expiry in seconds (2 hours)

	// Authentication
	JWT_EXPIRY_SECONDS = 86400 // JWT token expiry (24 hours)
)

// ================================
// NODE CODE MAPPINGS
// ================================

// NodeCodes maps QR code content to node numbers
var NodeCodes = map[string]int{
	// Generic format
	"NODE_001": 1,
	"NODE_002": 2,
	"NODE_003": 3,
	"NODE_004": 4,
	"NODE_005": 5,
	"NODE_006": 6,
	"NODE_007": 7,
}

// ================================
// CATEGORY CONSTANTS
// ================================

const (
	FUN_CATEGORY_NAME = "Fun" // Name of the fun category
)

// ================================
// DEFAULT VALUES
// ================================

const (
	DEFAULT_NODE_START = 0 // Initial node value (no node started yet)
	DEFAULT_SCORE      = 0 // Initial score value
	DEFAULT_RANK       = 0 // Default rank value
)

// ================================
// API RESPONSE CONSTANTS
// ================================

const (
	// Message templates
	WELCOME_NODE_MESSAGE     = "🎪 Welcome to Node %d! Answer all %d questions to continue your journey."
	COMPLETION_MESSAGE       = "🎉 Congratulations! You've completed all %d nodes of the carnival!"
	CORRECT_ANSWER_MESSAGE   = "✅ Correct! %d questions remaining in this node."
	INCORRECT_ANSWER_MESSAGE = "❌ Incorrect. %d questions remaining in this node."
	NODE_COMPLETED_MESSAGE   = "🎪 Node %d completed! Move to the next location to continue."
)

// ================================
// DATABASE CONSTANTS
// ================================

const (
	DEFAULT_DB_PORT = 5432 // Default PostgreSQL port
	DEFAULT_PORT    = 8080 // Default server port
)
