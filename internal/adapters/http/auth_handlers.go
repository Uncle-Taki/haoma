package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"haoma/internal/config"
	"haoma/internal/domain/player"
	"haoma/internal/infrastructure/auth"
)

// SignupRequest represents the user registration request
type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50" example:"Rostam Dastan"`
	Email    string `json:"email" binding:"required,email" example:"rostam@haoma.dev"`
	Password string `json:"password" binding:"required,min=6" example:"cyber_guardian_2024"`
}

// SignupResponse represents the response after successful registration
type SignupResponse struct {
	ID      uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name    string    `json:"name" example:"Rostam Dastan"`
	Email   string    `json:"email" example:"rostam@haoma.dev"`
	Message string    `json:"message" example:"Welcome to Haoma's carnival! Your account has been created."`
}

// LoginRequest represents the user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"rostam@haoma.dev"`
	Password string `json:"password" binding:"required" example:"cyber_guardian_2024"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Player      PlayerInfo `json:"player"`
	AccessToken string     `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string     `json:"token_type" example:"Bearer"`
	ExpiresIn   int        `json:"expires_in" example:"86400"` // seconds
	Message     string     `json:"message" example:"Welcome back to the carnival!"`
}

// PlayerInfo represents basic player information (no sensitive data)
type PlayerInfo struct {
	ID    uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name  string    `json:"name" example:"Rostam Dastan"`
	Email string    `json:"email" example:"rostam@haoma.dev"`
}

// Signup godoc
// @Summary Register a new player for the carnival
// @Description Create a new player account with name, email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SignupRequest true "Player registration information"
// @Success 201 {object} SignupResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/signup [post]
func (h *CarnivalHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	newPlayer, err := player.NewPlayer(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create player"})
		return
	}

	if err := h.service.CreatePlayer(newPlayer); err != nil {
		if err.Error() == "player already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save player"})
		return
	}

	c.JSON(http.StatusCreated, SignupResponse{
		ID:      newPlayer.ID,
		Name:    newPlayer.Name,
		Email:   newPlayer.Email,
		Message: "🎪 Welcome to Haoma's carnival! Your account has been created.",
	})
}

// Login godoc
// @Summary Authenticate a player
// @Description Login with email and password to verify player credentials
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Player login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *CarnivalHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	p, err := h.service.GetPlayerByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !p.ValidatePassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	jwtService := auth.NewJWTService(getJWTSecret())
	accessToken, err := jwtService.GeneratePlayerToken(p.ID, p.Name, p.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Player: PlayerInfo{
			ID:    p.ID,
			Name:  p.Name,
			Email: p.Email,
		},
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   config.JWT_EXPIRY_SECONDS,
		Message:     "🎪 Welcome back to the carnival!",
	})
}

// GetProfile godoc
// @Summary Get authenticated player profile information
// @Description Retrieve player profile from JWT token
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} PlayerInfo
// @Failure 401 {object} map[string]interface{}
// @Router /auth/profile [get]
func (h *CarnivalHandler) GetProfile(c *gin.Context) {

	playerID, exists := c.Get("player_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Player not authenticated"})
		return
	}

	playerName, _ := c.Get("player_name")
	playerEmail, _ := c.Get("player_email")

	c.JSON(http.StatusOK, PlayerInfo{
		ID:    playerID.(uuid.UUID),
		Name:  playerName.(string),
		Email: playerEmail.(string),
	})
}
