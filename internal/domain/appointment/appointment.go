package appointment

import (
	"time"

	"github.com/google/uuid"
)

// Status represents the current state of an appointment.
type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

// Appointment represents a scheduled veterinary visit at a customer's address.
type Appointment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	PetName   string    `gorm:"not null" json:"pet_name"`
	PetType   string    `gorm:"not null" json:"pet_type"` // e.g., "dog", "cat", "bird"
	Date      time.Time `gorm:"not null;index" json:"date"`
	TimeSlot  string    `gorm:"not null" json:"time_slot"` // e.g., "09:00-10:00"
	Address   string    `gorm:"not null" json:"address"`
	Notes     string    `json:"notes"`
	Status    Status    `gorm:"not null;default:pending" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewAppointment creates a new appointment with default pending status.
func NewAppointment(userID uuid.UUID, petName, petType string, date time.Time, timeSlot, address, notes string) *Appointment {
	return &Appointment{
		ID:        uuid.New(),
		UserID:    userID,
		PetName:   petName,
		PetType:   petType,
		Date:      date,
		TimeSlot:  timeSlot,
		Address:   address,
		Notes:     notes,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Cancel marks the appointment as cancelled.
func (a *Appointment) Cancel() {
	a.Status = StatusCancelled
	a.UpdatedAt = time.Now()
}

// Confirm marks the appointment as confirmed.
func (a *Appointment) Confirm() {
	a.Status = StatusConfirmed
	a.UpdatedAt = time.Now()
}

// Complete marks the appointment as completed.
func (a *Appointment) Complete() {
	a.Status = StatusCompleted
	a.UpdatedAt = time.Now()
}

// IsCancellable returns true if the appointment can still be cancelled.
func (a *Appointment) IsCancellable() bool {
	return a.Status == StatusPending || a.Status == StatusConfirmed
}
