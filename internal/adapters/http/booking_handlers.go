package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"haoma/internal/application/services"
	"haoma/internal/domain/user"
	"haoma/internal/infrastructure/auth"
)

// VetHandler handles HTTP requests for the veterinary care service.
type VetHandler struct {
	service *services.VetService
}

// --- Request / Response DTOs ---

// VetSignupRequest represents a new user registration request.
type VetSignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100" example:"Jane Doe"`
	Email    string `json:"email" binding:"required,email" example:"jane@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"securepass123"`
	Phone    string `json:"phone" example:"+1-555-0100"`
	Address  string `json:"address" example:"123 Pet Lane, Dogtown"`
}

// VetLoginRequest represents a user login request.
type VetLoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"jane@example.com"`
	Password string `json:"password" binding:"required" example:"securepass123"`
}

// VetLoginResponse represents the response after successful login.
type VetLoginResponse struct {
	User        VetUserInfo `json:"user"`
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type" example:"Bearer"`
	ExpiresIn   int         `json:"expires_in" example:"86400"`
}

// VetUserInfo contains non-sensitive user information.
type VetUserInfo struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	Address string    `json:"address"`
}

// BookAppointmentRequest represents a request to book a vet appointment.
type BookAppointmentRequest struct {
	PetName  string `json:"pet_name" binding:"required" example:"Buddy"`
	PetType  string `json:"pet_type" binding:"required" example:"dog"`
	Date     string `json:"date" binding:"required" example:"2026-04-15"`
	TimeSlot string `json:"time_slot" binding:"required" example:"10:00-11:00"`
	Address  string `json:"address" binding:"required" example:"456 Oak Avenue"`
	Notes    string `json:"notes" example:"Annual vaccination"`
}

// --- Handlers ---

// VetSignup registers a new user for the veterinary service.
// @Summary Register a new user
// @Description Create a new account for the veterinary care service
// @Tags Vet Auth
// @Accept json
// @Produce json
// @Param request body VetSignupRequest true "User registration info"
// @Success 201 {object} VetUserInfo
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /vet/auth/signup [post]
func (h *VetHandler) VetSignup(c *gin.Context) {
	var req VetSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	newUser, err := user.NewUser(req.Name, req.Email, req.Password, req.Phone, req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if err := h.service.CreateUser(newUser); err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, VetUserInfo{
		ID:      newUser.ID,
		Name:    newUser.Name,
		Email:   newUser.Email,
		Phone:   newUser.Phone,
		Address: newUser.Address,
	})
}

// VetLogin authenticates a user and returns a JWT token.
// @Summary Authenticate a user
// @Description Login with email and password to receive a JWT token
// @Tags Vet Auth
// @Accept json
// @Produce json
// @Param request body VetLoginRequest true "Login credentials"
// @Success 200 {object} VetLoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /vet/auth/login [post]
func (h *VetHandler) VetLogin(c *gin.Context) {
	var req VetLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	u, err := h.service.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !u.ValidatePassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	jwtService := auth.NewJWTService(getJWTSecret())
	accessToken, err := jwtService.GeneratePlayerToken(u.ID, u.Name, u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, VetLoginResponse{
		User: VetUserInfo{
			ID:      u.ID,
			Name:    u.Name,
			Email:   u.Email,
			Phone:   u.Phone,
			Address: u.Address,
		},
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   86400,
	})
}

// VetProfile returns the authenticated user's profile.
// @Summary Get user profile
// @Description Retrieve the authenticated user's profile information
// @Tags Vet Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} VetUserInfo
// @Failure 401 {object} map[string]interface{}
// @Router /vet/auth/profile [get]
func (h *VetHandler) VetProfile(c *gin.Context) {
	userID, exists := c.Get("player_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	u, err := h.service.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, VetUserInfo{
		ID:      u.ID,
		Name:    u.Name,
		Email:   u.Email,
		Phone:   u.Phone,
		Address: u.Address,
	})
}

// BookAppointment creates a new vet appointment.
// @Summary Book a vet appointment
// @Description Schedule a veterinary visit at your address
// @Tags Appointments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body BookAppointmentRequest true "Appointment details"
// @Success 201 {object} appointment.Appointment
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /vet/appointments [post]
func (h *VetHandler) BookAppointment(c *gin.Context) {
	userID, exists := c.Get("player_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	appt, err := h.service.BookAppointment(userID.(uuid.UUID), req.PetName, req.PetType, date, req.TimeSlot, req.Address, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appt)
}

// GetMyAppointments returns all appointments for the authenticated user.
// @Summary Get my appointments
// @Description Retrieve all appointments for the authenticated user
// @Tags Appointments
// @Security BearerAuth
// @Produce json
// @Success 200 {array} appointment.Appointment
// @Failure 401 {object} map[string]interface{}
// @Router /vet/appointments [get]
func (h *VetHandler) GetMyAppointments(c *gin.Context) {
	userID, exists := c.Get("player_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	appointments, err := h.service.GetUserAppointments(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}

// GetTodayAppointments returns all appointments for today (admin/vet dashboard).
// @Summary Get today's appointments
// @Description Retrieve all appointments scheduled for today (admin view)
// @Tags Appointments
// @Security BearerAuth
// @Produce json
// @Success 200 {array} appointment.Appointment
// @Router /vet/appointments/today [get]
func (h *VetHandler) GetTodayAppointments(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)
	appointments, err := h.service.GetAppointmentsByDate(today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}

// CancelAppointment cancels an existing appointment.
// @Summary Cancel an appointment
// @Description Cancel a pending or confirmed appointment
// @Tags Appointments
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /vet/appointments/:id/cancel [post]
func (h *VetHandler) CancelAppointment(c *gin.Context) {
	idParam := c.Param("id")
	appointmentID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	if err := h.service.CancelAppointment(appointmentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled successfully"})
}
