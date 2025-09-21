package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"haoma/internal/application/services"
	"haoma/internal/config"
	"haoma/internal/infrastructure/auth"
	"haoma/internal/infrastructure/persistence"
)

type CarnivalHandler struct {
	service *services.CarnivalService
}

func RegisterRoutes(router *gin.Engine, db *persistence.Database) {
	// Initialize repositories
	sessionRepo := persistence.NewSessionRepository(db.DB)
	questionRepo := persistence.NewQuestionRepository(db.DB)
	playerRepo := persistence.NewPlayerRepository(db.DB)
	leaderboardRepo := persistence.NewLeaderboardRepository(db.DB)

	// Initialize service
	service := services.NewCarnivalService(sessionRepo, questionRepo, playerRepo, leaderboardRepo)

	// Initialize handler
	handler := &CarnivalHandler{service: service}

	// Initialize JWT service and middleware
	jwtService := auth.NewJWTService(getJWTSecret())
	jwtMiddleware := auth.JWTMiddleware(jwtService)

	// API routes
	api := router.Group("/api/v1")
	{
		// Public authentication routes (no JWT required)
		authPublic := api.Group("/auth")
		{
			authPublic.POST("/signup", handler.Signup)
			authPublic.POST("/login", handler.Login)
		}

		// Protected authentication routes (JWT required)
		authProtected := api.Group("/auth")
		authProtected.Use(jwtMiddleware)
		{
			authProtected.GET("/profile", handler.GetProfile)
		}

		// Protected game session routes (JWT required)
		sessions := api.Group("/sessions")
		sessions.Use(jwtMiddleware)
		{
			sessions.POST("/start", handler.StartSession)
			sessions.POST("/:id/answer", handler.SubmitAnswer)
		}

		// Protected node access via QR codes (JWT required)
		nodes := api.Group("/nodes")
		nodes.Use(jwtMiddleware)
		{
			nodes.POST("/scan", handler.ScanNodeQR)
		}

		// Public leaderboard (no authentication needed)
		api.GET("/leaderboard", handler.GetLeaderboard)
	}
}

// StartSessionRequest represents the request to begin a carnival journey (no player_id needed - from JWT)
type StartSessionRequest struct {
}

// StartNodeRequest represents scanning a QR code to start a specific node (no player_id needed - from JWT)
type StartNodeRequest struct {
	NodeCode  string     `json:"node_code" binding:"required" example:"NODE_CRYPTO_001"`
	SessionID *uuid.UUID `json:"session_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"`
}

// StartSessionResponse represents the response when starting a session
type StartSessionResponse struct {
	SessionID uuid.UUID `json:"session_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Message   string    `json:"message" example:"Session created! Scan a node QR code to begin your journey."`
}

// StartNodeResponse represents the response when starting a node via QR code
type StartNodeResponse struct {
	SessionID uuid.UUID    `json:"session_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	Node      NodeResponse `json:"node"`
	Message   string       `json:"message" example:"🎪 Welcome to Node 1: Cryptography! Answer all questions to continue."`
}

// NodeResponse represents a carnival node (tent) with questions
type NodeResponse struct {
	Number              int                `json:"number" example:"1"`
	CategoryName        string             `json:"category_name" example:"Cryptography"`
	CategoryDescription string             `json:"category_description" example:"Learn about encryption, decryption, and cryptographic protocols"`
	Questions           []QuestionResponse `json:"questions"`
}

// QuestionResponse represents a question without the correct answer
type QuestionResponse struct {
	ID      uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Text    string    `json:"text" example:"What does SQL injection exploit?"`
	OptionA string    `json:"option_a" example:"Input validation"`
	OptionB string    `json:"option_b" example:"Database queries"`
	OptionC *string   `json:"option_c,omitempty" example:"File uploads"`
	OptionD *string   `json:"option_d,omitempty" example:"Network protocols"`
}

// StartSession godoc
// @Summary Start a new carnival session
// @Description Begin a mystical journey through carnival nodes of cyber trials
// @Tags Sessions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body StartSessionRequest true "Player information"
// @Success 200 {object} StartSessionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sessions/start [post]
func (h *CarnivalHandler) StartSession(c *gin.Context) {
	playerID, exists := c.Get("player_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Player not authenticated"})
		return
	}

	session, err := h.service.CreateSession(playerID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, StartSessionResponse{
		SessionID: session.ID,
		Message:   "🎪 Session created! Scan a node QR code at any carnival location to begin your journey.",
	})
}

// SubmitAnswerRequest represents an answer submission
type SubmitAnswerRequest struct {
	QuestionID uuid.UUID `json:"question_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	Answer     string    `json:"answer" binding:"required" example:"B"`
}

// SubmitAnswerResponse represents the response after answering
type SubmitAnswerResponse struct {
	IsCorrect     bool   `json:"is_correct" example:"true"`
	NodeCompleted bool   `json:"node_completed" example:"false"`
	Message       string `json:"message" example:"Correct! 4 questions remaining in this node."`
	CurrentScore  *int   `json:"current_score,omitempty" example:"320"`
}

// SubmitAnswer godoc
// @Summary Submit an answer to a question
// @Description Answer a riddle from the current carnival node
// @Tags Sessions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body SubmitAnswerRequest true "Answer submission"
// @Success 200 {object} SubmitAnswerResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sessions/{id}/answer [post]
func (h *CarnivalHandler) SubmitAnswer(c *gin.Context) {
	sessionIDParam := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.SubmitAnswer(sessionID, req.QuestionID, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var message string
	if result.NodeCompleted {
		message = "🎪 Node completed! Check your updated leaderboard position. Find the next location to continue."
	} else {
		remaining := config.QUESTIONS_PER_NODE - result.QuestionsAnsweredInNode
		if result.IsCorrect {
			message = fmt.Sprintf(config.CORRECT_ANSWER_MESSAGE, remaining)
		} else {
			message = fmt.Sprintf(config.INCORRECT_ANSWER_MESSAGE, remaining)
		}
	}

	response := SubmitAnswerResponse{
		IsCorrect:     result.IsCorrect,
		NodeCompleted: result.NodeCompleted,
		Message:       message,
	}

	if result.NodeCompleted {
		response.CurrentScore = &result.CurrentScore
	}

	c.JSON(http.StatusOK, response)
}

// LeaderboardResponse represents the taxteh-ye sharaf
type LeaderboardResponse struct {
	Entries []LeaderboardEntry `json:"entries"`
}

// LeaderboardEntry represents a champion's achievement
type LeaderboardEntry struct {
	Rank           int    `json:"rank" example:"1"`
	PlayerName     string `json:"player_name" example:"Rostam"`
	FinalScore     int    `json:"final_score" example:"850"`
	CompletionTime string `json:"completion_time" example:"38m45s"`
	AchievedAt     string `json:"achieved_at" example:"2025-09-18T14:30:45Z"`
}

// GetLeaderboard godoc
// @Summary Get the top leaderboard
// @Description Retrieve the taxteh-ye sharaf showing the greatest champions
// @Tags Leaderboard
// @Produce json
// @Success 200 {object} LeaderboardResponse
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard [get]
func (h *CarnivalHandler) GetLeaderboard(c *gin.Context) {
	entries, err := h.service.GetLeaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := LeaderboardResponse{
		Entries: make([]LeaderboardEntry, len(entries)),
	}

	for i, entry := range entries {
		resp.Entries[i] = LeaderboardEntry{
			Rank:           i + config.MIN_NODE_NUMBER,
			PlayerName:     entry.PlayerName,
			FinalScore:     entry.FinalScore,
			CompletionTime: entry.CompletionTime.String(),
			AchievedAt:     entry.AchievedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, resp)
}
